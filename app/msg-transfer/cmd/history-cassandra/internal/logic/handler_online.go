package logic

import (
	"context"
	"fmt"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-cassandra/internal/repository"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-cassandra/internal/svc"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/Path-IM/Path-IM-Server/common/utils/statistics"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"
	"github.com/golang/protobuf/proto"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	singleMsgSuccessCount uint64
	groupMsgCount         uint64
	singleMsgFailedCount  uint64
)

func init() {
	statistics.NewStatistics(&singleMsgSuccessCount, "msg-transfer-history-cassandra", fmt.Sprintf("%d second singleMsgCount insert to cassandra", 300), 300)
	statistics.NewStatistics(&groupMsgCount, "msg-transfer-history-cassandra", fmt.Sprintf("%d second groupMsgCount insert to cassandra", 300), 300)
	statistics.NewStatistics(&groupMsgCount, "msg-transfer-history-cassandra", fmt.Sprintf("%d second groupMsgCount insert to cassandra", 300), 300)
}

type MsgTransferHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	rep *repository.Rep
}

func NewMsgTransferHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MsgTransferHistoryLogic {
	return &MsgTransferHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		rep:    repository.NewRep(svcCtx),
	}
}

func (l *MsgTransferHistoryLogic) saveUserChat(ctx context.Context, uid string, msg *chatpb.MsgDataToMQ) error {
	var seq uint64
	var err error
	xtrace.StartFuncSpan(ctx, "MsgTransferHistoryLogic.saveUserChat.IncrUserSeq", func(ctx context.Context) {
		seq, err = l.rep.IncrUserSeq(uid)
	})
	if err != nil {
		l.Logger.Error("data insert to redis err ", err.Error(), msg.String())
		return err
	}
	msg.MsgData.Seq = uint32(seq)
	pbSaveData := chatpb.MsgDataToDB{}
	pbSaveData.MsgData = msg.MsgData
	xtrace.StartFuncSpan(ctx, "MsgTransferHistoryLogic.saveUserChat.SaveUserChatCassandra2", func(ctx context.Context) {
		err = l.rep.SaveUserChatCassandra2(ctx, uid, int64(pbSaveData.MsgData.ServerTime), &pbSaveData)
	})
	return err
}
func (l *MsgTransferHistoryLogic) saveGroupChat(ctx context.Context, groupId string, msg *chatpb.MsgDataToMQ) error {
	var seq uint64
	var err error
	xtrace.StartFuncSpan(ctx, "MsgTransferHistoryLogic.saveGroupChat.IncrUserSeq", func(ctx context.Context) {
		seq, err = l.rep.IncrGroupSeq(groupId)
	})
	if err != nil {
		l.Logger.Error("data insert to redis err ", err.Error(), msg.String())
		return err
	}
	msg.MsgData.Seq = uint32(seq)
	pbSaveData := chatpb.MsgDataToDB{}
	pbSaveData.MsgData = msg.MsgData
	xtrace.StartFuncSpan(ctx, "MsgTransferHistoryLogic.saveGroupChat.SaveGroupChatCassandra", func(ctx context.Context) {
		err = l.rep.SaveGroupChatCassandra(ctx, groupId, int64(pbSaveData.MsgData.ServerTime), &pbSaveData)
	})
	return err
}

func (l *MsgTransferHistoryLogic) ChatMs2Cassandra(msg []byte, msgKey string) (err error) {
	msgFromMQ := chatpb.MsgDataToMQ{}
	xtrace.StartFuncSpan(l.ctx, "MsgTransferHistoryLogic.ChatMs2Cassandra.UnmarshalMsg", func(ctx context.Context) {
		err = proto.Unmarshal(msg, &msgFromMQ)
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("unmarshal msg failed, err: %v", err)
		return
	}
	logx.WithContext(l.ctx).Infof("msg: %v", msgFromMQ.String())
	isHistory := chatpb.GetSwitchFromOptions(msgFromMQ.MsgData.MsgOptions, types.IsHistory, true)
	isSenderSync := chatpb.GetSwitchFromOptions(msgFromMQ.MsgData.MsgOptions, types.IsSenderSync, true)
	switch msgFromMQ.MsgData.ConversationType {
	case types.SingleChatType:
		xtrace.StartFuncSpan(l.ctx, "MsgTransferHistoryLogic.ChatMs2Cassandra.SingleChat", func(ctx context.Context) {
			if isHistory {
				err = l.saveUserChat(ctx, msgKey, &msgFromMQ)
				if err != nil {
					singleMsgFailedCount++
					l.Logger.Error("single data insert to cassandra err ", err.Error(), " ", msgFromMQ.String())
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
		xtrace.StartFuncSpan(l.ctx, "MsgTransferHistoryLogic.ChatMs2Cassandra.GroupChat", func(ctx context.Context) {
			if isHistory {
				err = l.saveGroupChat(ctx, msgFromMQ.MsgData.ReceiveID, &msgFromMQ)
				if err != nil {
					l.Logger.Error("group data insert to cassandra err ", msgFromMQ.String(), " GroupID ", msgFromMQ.MsgData.ReceiveID, " ", err.Error())
					return
				}
				groupMsgCount++
			}
			go l.sendMessageToGroupPush(ctx, &msgFromMQ, msgFromMQ.MsgData.ReceiveID)
		})
	default:
		l.Logger.Error("SessionType error ", msgFromMQ.String())
		return
	}
	return err
}
