package logic

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	pushpb "github.com/showurl/Zero-IM-Server/app/msg-push/cmd/rpc/pb"
	"github.com/showurl/Zero-IM-Server/app/msg-transfer/cmd/history/internal/repository"
	"github.com/showurl/Zero-IM-Server/app/msg-transfer/cmd/history/internal/svc"
	chatpb "github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/showurl/Zero-IM-Server/common/types"
	"github.com/showurl/Zero-IM-Server/common/utils"
	"github.com/showurl/Zero-IM-Server/common/utils/statistics"
	"github.com/showurl/Zero-IM-Server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	singleMsgSuccessCount uint64
	groupMsgCount         uint64
	superGroupMsgCount    uint64
	singleMsgFailedCount  uint64
)

func init() {
	statistics.NewStatistics(&singleMsgSuccessCount, "msg-transfer-history", fmt.Sprintf("%d second singleMsgCount insert to mongo", 300), 300)
	statistics.NewStatistics(&groupMsgCount, "msg-transfer-history", fmt.Sprintf("%d second groupMsgCount insert to mongo", 300), 300)
	statistics.NewStatistics(&superGroupMsgCount, "msg-transfer-history", fmt.Sprintf("%d second superGroupMsgCount insert to mongo", 300), 300)
}

type MsgTransferHistoryOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	rep *repository.Rep
}

func NewMsgTransferHistoryOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MsgTransferHistoryOnlineLogic {
	return &MsgTransferHistoryOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		rep:    repository.NewRep(svcCtx),
	}
}

func (l *MsgTransferHistoryOnlineLogic) saveUserChat(ctx context.Context, uid string, msg *chatpb.MsgDataToMQ) error {
	var seq uint64
	var err error
	xtrace.StartFuncSpan(ctx, "MsgTransferHistoryOnlineLogic.saveUserChat.IncrUserSeq", func(ctx context.Context) {
		seq, err = l.rep.IncrUserSeq(uid)
	})
	if err != nil {
		l.Logger.Error("data insert to redis err ", err.Error(), msg.String())
		return err
	}
	msg.MsgData.Seq = uint32(seq)
	pbSaveData := chatpb.MsgDataToDB{}
	pbSaveData.MsgData = msg.MsgData
	xtrace.StartFuncSpan(ctx, "MsgTransferHistoryOnlineLogic.saveUserChat.SaveUserChatMongo2", func(ctx context.Context) {
		err = l.rep.SaveUserChatMongo2(ctx, uid, pbSaveData.MsgData.SendTime, &pbSaveData)
	})
	return err
}
func (l *MsgTransferHistoryOnlineLogic) saveSuperGroupChat(ctx context.Context, groupId string, msg *chatpb.MsgDataToMQ) error {
	var seq uint64
	var err error
	xtrace.StartFuncSpan(ctx, "MsgTransferHistoryOnlineLogic.saveSuperGroupChat.IncrUserSeq", func(ctx context.Context) {
		seq, err = l.rep.IncrSuperGroupSeq(groupId)
	})
	if err != nil {
		l.Logger.Error("data insert to redis err ", err.Error(), msg.String())
		return err
	}
	msg.MsgData.Seq = uint32(seq)
	pbSaveData := chatpb.MsgDataToDB{}
	pbSaveData.MsgData = msg.MsgData
	xtrace.StartFuncSpan(ctx, "MsgTransferHistoryOnlineLogic.saveSuperGroupChat.SaveSuperGroupChatMongo2", func(ctx context.Context) {
		err = l.rep.SaveSuperGroupChatMongo2(ctx, groupId, pbSaveData.MsgData.SendTime, &pbSaveData)
	})
	return err
}
func (l *MsgTransferHistoryOnlineLogic) sendMessageToPush(ctx context.Context, message *chatpb.MsgDataToMQ, pushToUserID string) {
	logx.WithContext(ctx).Info("msg_transfer send message to push", "message", message.String())
	rpcPushMsg := pushpb.PushMsgReq{MsgData: message.MsgData, PushToUserID: pushToUserID}
	_, err := l.svcCtx.MsgPush.PushMsg(ctx, &rpcPushMsg)
	if err != nil {
		logx.WithContext(ctx).Error("rpc send failed", "push data", rpcPushMsg.String(), "err", err.Error())
		mqPushMsg := chatpb.PushMsgDataToMQ{MsgData: message.MsgData, PushToUserID: pushToUserID, OperationID: xtrace.TraceIdFromContext(l.ctx)}
		pid, offset, err := l.svcCtx.SinglePushProducer.SendMessage(ctx, &mqPushMsg)
		if err != nil {
			logx.WithContext(ctx).Error("kafka send failed", mqPushMsg.OperationID, "send data", mqPushMsg.String(), "pid", pid, "offset", offset, "err", err.Error())
		}
	} else {
		logx.WithContext(ctx).Info("rpc send success", "push data", rpcPushMsg.String())
	}
}

func (l *MsgTransferHistoryOnlineLogic) sendMessageToSuperGroupPush(ctx context.Context, message *chatpb.MsgDataToMQ, groupId string) {
	mqPushMsg := chatpb.PushMsgToSuperGroupDataToMQ{MsgData: message.MsgData, SuperGroupID: groupId, OperationID: xtrace.TraceIdFromContext(l.ctx)}
	pid, offset, err := l.svcCtx.SuperGroupPushProducer.SendMessage(ctx, &mqPushMsg)
	if err != nil {
		logx.WithContext(ctx).Error("kafka send failed ", "send data ", mqPushMsg.String(), " pid ", pid, " offset ", offset, " err ", err.Error())
	}
}

func (l *MsgTransferHistoryOnlineLogic) ChatMs2Mongo(msg []byte, msgKey string) {
	msgFromMQ := chatpb.MsgDataToMQ{}
	var err error
	xtrace.StartFuncSpan(l.ctx, "MsgTransferHistoryOnlineLogic.ChatMs2Mongo.UnmarshalMsg", func(ctx context.Context) {
		err = proto.Unmarshal(msg, &msgFromMQ)
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("unmarshal msg failed, err: %v", err)
		return
	}
	logx.WithContext(l.ctx).Infof("msg: %v", msgFromMQ.String())
	isHistory := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, types.IsHistory)
	isSenderSync := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, types.IsSenderSync)
	switch msgFromMQ.MsgData.SessionType {
	case types.SingleChatType:
		xtrace.StartFuncSpan(l.ctx, "MsgTransferHistoryOnlineLogic.ChatMs2Mongo.SingleChat", func(ctx context.Context) {
			if isHistory {
				err := l.saveUserChat(ctx, msgKey, &msgFromMQ)
				if err != nil {
					singleMsgFailedCount++
					l.Logger.Error("single data insert to mongo err ", err.Error(), " ", msgFromMQ.String())
					return
				}
				singleMsgSuccessCount++
			}
			if !isSenderSync && msgKey == msgFromMQ.MsgData.SendID {
			} else {
				go l.sendMessageToPush(ctx, &msgFromMQ, msgKey)
			}
		})
	case types.GroupChatType:
		xtrace.StartFuncSpan(l.ctx, "MsgTransferHistoryOnlineLogic.ChatMs2Mongo.SuperGroupChat", func(ctx context.Context) {
			if isHistory {
				err := l.saveSuperGroupChat(ctx, msgFromMQ.MsgData.GroupID, &msgFromMQ)
				if err != nil {
					l.Logger.Error("super group data insert to mongo err ", msgFromMQ.String(), " GroupID ", msgFromMQ.MsgData.GroupID, " ", err.Error())
					return
				}
				superGroupMsgCount++
			}
			go l.sendMessageToSuperGroupPush(ctx, &msgFromMQ, msgFromMQ.MsgData.GroupID)
		})
	case types.NotificationChatType:
		xtrace.StartFuncSpan(l.ctx, "MsgTransferHistoryOnlineLogic.ChatMs2Mongo.NotificationChat", func(ctx context.Context) {
			if isHistory {
				err := l.saveUserChat(ctx, msgKey, &msgFromMQ)
				if err != nil {
					l.Logger.Error("single data insert to mongo err ", err.Error(), " ", msgFromMQ.String())
					return
				}
			}
			if !isSenderSync && msgKey == msgFromMQ.MsgData.SendID {
			} else {
				go l.sendMessageToPush(ctx, &msgFromMQ, msgKey)
			}
		})
	default:
		l.Logger.Error("SessionType error ", msgFromMQ.String())
		return
	}
}
