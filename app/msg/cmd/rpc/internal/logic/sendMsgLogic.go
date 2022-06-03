package logic

import (
	"context"
	imuserpb "github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/Path-IM/Path-IM-Server/common/utils"

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
	//if !utils.VerifyToken(pb.Token, pb.SendID) {
	//	return returnMsg(&replay, pb, http.StatusUnauthorized, "token validate err,not authorized", "", 0)
	l.encapsulateMsgData(pb.MsgData)
	logx.WithContext(l.ctx).Info("this is a test MsgData ", pb.MsgData)
	msgToMQSingle := chatpb.MsgDataToMQ{Token: pb.Token, MsgData: pb.MsgData}
	//options := utils.JsonStringToMap(pbData.Options)
	isHistory := utils.GetSwitchFromOptions(pb.MsgData.Options, types.IsHistory)
	mReq := MsgCallBackReq{
		SendID:      pb.MsgData.SendID,
		RecvID:      pb.MsgData.RecvID,
		Content:     string(pb.MsgData.Content),
		SendTime:    pb.MsgData.ServerTime,
		MsgFrom:     pb.MsgData.MsgFrom,
		ContentType: pb.MsgData.ContentType,
		SessionType: pb.MsgData.SessionType,
		PlatformID:  pb.MsgData.SenderPlatformID,
		MsgID:       pb.MsgData.ClientMsgID,
	}
	if !isHistory {
		mReq.IsOnlineOnly = true
	}

	// callback
	canSend, err := l.callbackWordFilter(pb)
	if err != nil {
		logx.WithContext(l.ctx).Error(utils.GetSelfFuncName(), "callbackWordFilter failed", err.Error(), pb.MsgData)
	}
	if !canSend {
		return returnMsg(&replay, pb, 201, "callbackWordFilter result stop rpc and return", "", 0)
	}
	switch pb.MsgData.SessionType {
	case types.SingleChatType:
		// callback
		canSend, err := l.callbackBeforeSendSingleMsg(pb)
		if err != nil {
			logx.WithContext(l.ctx).Error(utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg failed", err.Error())
		}
		if !canSend {
			return returnMsg(&replay, pb, 201, "callbackBeforeSendSingleMsg result stop rpc and return", "", 0)
		}
		isSend := l.modifyMessageByUserMessageReceiveOpt(pb.MsgData.RecvID, pb.MsgData.SendID, types.SingleChatType, pb)
		if isSend {
			msgToMQSingle.MsgData = pb.MsgData
			logx.WithContext(l.ctx).Info(msgToMQSingle.String())
			err1 := l.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.RecvID, types.OnlineStatus)
			if err1 != nil {
				logx.WithContext(l.ctx).Error(msgToMQSingle.OperationID, "kafka send msg err:RecvID ", msgToMQSingle.MsgData.RecvID, msgToMQSingle.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err ", "", 0)
			}
		}
		if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.RecvID { //Filter messages sent to yourself
			err2 := l.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.SendID, types.OnlineStatus)
			if err2 != nil {
				logx.WithContext(l.ctx).Error(msgToMQSingle.OperationID, "kafka send msg err:SendID ", msgToMQSingle.MsgData.SendID, msgToMQSingle.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err ", "", 0)
			}
		}
		// callback
		if err := l.callbackAfterSendSingleMsg(pb); err != nil {
			logx.WithContext(l.ctx).Error(utils.GetSelfFuncName(), "callbackAfterSendSingleMsg failed", err.Error())
		}
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.ServerTime)

	case types.GroupChatType:
		// callback
		canSend, err := l.callbackBeforeSendSuperGroupMsg(pb)
		if err != nil {
			logx.WithContext(l.ctx).Error(utils.GetSelfFuncName(), "callbackBeforeSendSuperGroupMsg failed ", err.Error())
		}
		if !canSend {
			return returnMsg(&replay, pb, 201, "callbackBeforeSendSuperGroupMsg result stop rpc and return", " ", 0)
		}
		// 读扩散
		msgToMQSingle.MsgData = pb.MsgData
		logx.WithContext(l.ctx).Info(msgToMQSingle.String())
		err1 := l.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.GroupID, types.OnlineStatus)
		if err1 != nil {
			logx.WithContext(l.ctx).Error(msgToMQSingle.OperationID, "kafka send msg err:GroupID ", msgToMQSingle.MsgData.GroupID, msgToMQSingle.String())
			return returnMsg(&replay, pb, 201, "kafka send msg err ", "", 0)
		}
		// callback
		if err := l.callbackAfterSendSuperGroupMsg(pb); err != nil {
			logx.WithContext(l.ctx).Error(utils.GetSelfFuncName(), "callbackAfterSendSuperGroupMsg failed ", err.Error())
		}
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.ServerTime)
	case types.NotificationChatType:
		msgToMQSingle.MsgData = pb.MsgData
		logx.WithContext(l.ctx).Info(msgToMQSingle.OperationID, msgToMQSingle.String())
		err1 := l.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.RecvID, types.OnlineStatus)
		if err1 != nil {
			logx.WithContext(l.ctx).Error(msgToMQSingle.OperationID, "kafka send msg err:RecvID", msgToMQSingle.MsgData.RecvID, msgToMQSingle.String())
			return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
		}

		if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.RecvID { //Filter messages sent to yourself
			err2 := l.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.SendID, types.OnlineStatus)
			if err2 != nil {
				logx.WithContext(l.ctx).Error(msgToMQSingle.OperationID, "kafka send msg err:SendID", msgToMQSingle.MsgData.SendID, msgToMQSingle.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.ServerTime)
	default:
		return returnMsg(&replay, pb, 203, "unkonwn sessionType", "", 0)
	}
}

func (l *SendMsgLogic) getOnlineAndOfflineUserIDList(list []string) (online []string, offline []string) {
	return list, nil
}

func returnMsg(replay *chatpb.SendMsgResp, pb *chatpb.SendMsgReq, errCode int32, errMsg, serverMsgID string, sendTime int64) (*chatpb.SendMsgResp, error) {
	replay.ErrCode = errCode
	replay.ErrMsg = errMsg
	replay.ServerMsgID = serverMsgID
	replay.ClientMsgID = pb.MsgData.ClientMsgID
	replay.ServerTime = sendTime
	replay.SendID = pb.MsgData.SendID
	replay.RecvID = pb.MsgData.RecvID
	replay.GroupID = pb.MsgData.GroupID
	return replay, nil
}

func (l *SendMsgLogic) userRelationshipVerification(data *chatpb.SendMsgReq) (bool, int32, string) {
	if data.MsgData.SessionType == types.GroupChatType {
		return true, 0, ""
	}
	// 是不是拉黑了
	ifInBlackResp, err := l.svcCtx.ImUser.IfAInBBlacklist(l.ctx, &imuserpb.IfAInBBlacklistReq{
		AUserID: data.MsgData.SendID,
		BUserID: data.MsgData.RecvID,
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
		needFriend := true
		switch data.MsgData.ContentType {
		case types.NotificationUser2User:
			needFriend = false
		}
		if !needFriend {
			return true, 0, ""
		}
		// 是不是好友
		ifInFriendResp, err := l.svcCtx.ImUser.IfAInBFriendList(l.ctx, &imuserpb.IfAInBFriendListReq{
			AUserID: data.MsgData.SendID,
			BUserID: data.MsgData.RecvID,
		})
		if err != nil {
			logx.WithContext(l.ctx).Error("GetFriendIDListFromCache rpc call failed ", err.Error())
		} else {
			if ifInFriendResp.CommonResp.ErrCode != 0 {
				logx.WithContext(l.ctx).Error("GetFriendIDListFromCache rpc logic call failed ", ifInFriendResp.String())
			} else {
				if !ifInFriendResp.IsInFriendList {
					return false, 601, "not friend"
				}
			}
		}
		return true, 0, ""
	} else {
		return true, 0, ""
	}
}
