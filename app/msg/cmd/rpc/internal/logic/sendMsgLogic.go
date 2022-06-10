package logic

import (
	"context"
	imuserpb "github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/Path-IM/Path-IM-Server/common/utils"
	timeUtils "github.com/Path-IM/Path-IM-Server/common/utils/time"
	"github.com/Path-IM/Path-IM-Server/common/xorm/global"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMsgLogic) SendMsg(pb *pb.SendMsgReq) (*pb.SendMsgResp, error) {
	replay := chatpb.SendMsgResp{}
	flag, errCode, errMsg := l.userRelationshipVerification(pb)
	if !flag {
		return returnMsg(&replay, pb, errCode, errMsg, "", 0)
	}
	// 消息预处理
	{
		// 生成消息id
		pb.MsgData.ServerMsgID = global.GetID()
		pb.MsgData.ServerTime = timeUtils.GetCurrentTimestampByMill()
		// 修改消息options
		l.encapsulateMsgData(pb.MsgData)
		logx.WithContext(l.ctx).Info("this is a test MsgData ", pb.MsgData)
	}
	msgToMQSingle := chatpb.MsgDataToMQ{MsgData: pb.MsgData}
	switch pb.MsgData.ConversationType {
	case types.SingleChatType:
		// callback
		{
			canSend, err := l.callbackBeforeSendSingleMsg(pb)
			if err != nil {
				logx.WithContext(l.ctx).Error(utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg failed", err.Error())
			}
			if !canSend {
				return returnMsg(&replay, pb, types.ErrCodeFailed, "callbackBeforeSendSingleMsg result stop rpc and return", "", 0)
			}
		}
		isSend := l.modifyMessageByUserMessageReceiveOpt(pb.MsgData.ReceiveID, pb.MsgData.SendID, types.SingleChatType, pb)
		if isSend {
			msgToMQSingle.MsgData = pb.MsgData
			logx.WithContext(l.ctx).Info(msgToMQSingle.String())
			err1 := l.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.ReceiveID)
			if err1 != nil {
				logx.WithContext(l.ctx).Error(msgToMQSingle.TraceId, "kafka send msg err:RecvID ", msgToMQSingle.MsgData.ReceiveID, msgToMQSingle.String())
				return returnMsg(&replay, pb, types.ErrCodeFailed, "kafka send msg err ", "", 0)
			}
		}
		if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.ReceiveID { //Filter messages sent to yourself
			err2 := l.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.SendID)
			if err2 != nil {
				logx.WithContext(l.ctx).Error(msgToMQSingle.TraceId, "kafka send msg err:SendID ", msgToMQSingle.MsgData.SendID, msgToMQSingle.String())
				return returnMsg(&replay, pb, types.ErrCodeFailed, "kafka send msg err ", "", 0)
			}
		}
		// callback
		if err := l.callbackAfterSendSingleMsg(pb); err != nil {
			logx.WithContext(l.ctx).Error(utils.GetSelfFuncName(), "callbackAfterSendSingleMsg failed", err.Error())
		}
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, int64(msgToMQSingle.MsgData.ServerTime))

	case types.GroupChatType:
		// callback
		{
			canSend, err := l.callbackBeforeSendGroupMsg(pb)
			if err != nil {
				logx.WithContext(l.ctx).Error(utils.GetSelfFuncName(), "callbackBeforeSendGroupMsg failed ", err.Error())
			}
			if !canSend {
				return returnMsg(&replay, pb, types.ErrCodeFailed, "callbackBeforeSendGroupMsg result stop rpc and return", " ", 0)
			}
		}
		// 读扩散
		msgToMQSingle.MsgData = pb.MsgData
		logx.WithContext(l.ctx).Info(msgToMQSingle.String())
		err1 := l.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.ReceiveID)
		if err1 != nil {
			logx.WithContext(l.ctx).Error(msgToMQSingle.TraceId, " kafka send msg err:GroupID ", msgToMQSingle.MsgData.ReceiveID, msgToMQSingle.String())
			return returnMsg(&replay, pb, types.ErrCodeFailed, "kafka send msg err", "", 0)
		}
		// callback
		if err := l.callbackAfterSendGroupMsg(pb); err != nil {
			logx.WithContext(l.ctx).Error(utils.GetSelfFuncName(), "callbackAfterSendGroupMsg failed ", err.Error())
		}
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, int64(msgToMQSingle.MsgData.ServerTime))
	default:
		return returnMsg(&replay, pb, types.ErrCodeFailed, "unkonwn sessionType", "", 0)
	}
}

func returnMsg(replay *chatpb.SendMsgResp, pb *chatpb.SendMsgReq, errCode int32, errMsg, serverMsgID string, sendTime int64) (*chatpb.SendMsgResp, error) {
	replay.ErrCode = errCode
	replay.ErrMsg = errMsg
	replay.ServerMsgID = serverMsgID
	replay.ClientMsgID = pb.MsgData.ClientMsgID
	replay.ServerTime = sendTime
	replay.ReceiveID = pb.MsgData.ReceiveID
	replay.ContentType = pb.MsgData.ContentType
	replay.ConversationType = pb.MsgData.ConversationType
	return replay, nil
}

func (l *SendMsgLogic) userRelationshipVerification(data *chatpb.SendMsgReq) (bool, int32, string) {
	if data.MsgData.ConversationType == types.GroupChatType {
		return true, 0, ""
	}
	// 是不是拉黑了
	ifInBlackResp, err := l.svcCtx.ImUser.IfAInBBlacklist(l.ctx, &imuserpb.IfAInBBlacklistReq{
		AUserID: data.MsgData.SendID,
		BUserID: data.MsgData.ReceiveID,
	})
	if err != nil {
		logx.WithContext(l.ctx).Error("GetBlackIDListFromCache rpc call failed ", err.Error())
	} else {
		if ifInBlackResp.CommonResp.ErrCode != 0 {
			logx.WithContext(l.ctx).Error("GetBlackIDListFromCache rpc logic call failed ", ifInBlackResp.String())
		} else {
			if ifInBlackResp.IsInBlacklist {
				return false, 600, "in black list"
			}
		}
	}
	if l.svcCtx.Config.MessageVerify.FriendVerify {
		needFriend := pb.GetSwitchFromOptions(data.MsgData.MsgOptions, types.NeedBeFriend, false)
		if !needFriend {
			return true, 0, ""
		}
		// 是不是好友
		ifInFriendResp, err := l.svcCtx.ImUser.IfAInBFriendList(l.ctx, &imuserpb.IfAInBFriendListReq{
			AUserID: data.MsgData.SendID,
			BUserID: data.MsgData.ReceiveID,
		})
		if err != nil {
			logx.WithContext(l.ctx).Error("GetFriendIDListFromCache rpc call failed ", err.Error())
		} else {
			if ifInFriendResp.CommonResp.ErrCode != 0 {
				logx.WithContext(l.ctx).Error("GetFriendIDListFromCache rpc logic call failed ", ifInFriendResp.String())
			} else {
				if !ifInFriendResp.IsInFriendList {
					return false, types.ErrCodeFailed, "not friend"
				}
			}
		}
		return true, 0, ""
	} else {
		return true, 0, ""
	}
}

func (l *SendMsgLogic) modifyMessageByUserMessageReceiveOpt(userID, sourceID string, sessionType int, pb *chatpb.SendMsgReq) bool {
	// 用户设置了消息接收选项
	req := &imuserpb.GetSingleConversationRecvMsgOptsReq{
		UserID:       userID,
		SenderUserID: sourceID,
	}
	resp, err := l.svcCtx.ImUser.GetSingleConversationRecvMsgOpts(l.ctx, req)
	if err != nil {
		logx.WithContext(l.ctx).Error("GetSingleConversationMsgOpt from redis err ", pb.String(), " ", err.Error())
		return true
	} else if resp.CommonResp == nil || resp.CommonResp.ErrCode != 0 {
		return true
	} else {
		switch resp.Opts {
		case imuserpb.RecvMsgOpt_NotReceiveMessage:
			return false
		case imuserpb.RecvMsgOpt_ReceiveNotNotifyMessage:
			if pb.MsgData.MsgOptions == nil {
				pb.MsgData.MsgOptions = &chatpb.MsgOptions{}
			}
			chatpb.SetSwitchFromOptions(pb.MsgData.MsgOptions, types.IsOfflinePush, false)
		}
		return true
	}
}
