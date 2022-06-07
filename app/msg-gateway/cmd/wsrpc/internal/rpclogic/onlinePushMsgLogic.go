package rpclogic

import (
	"context"
	imuserpb "github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcsvc"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnlinePushMsgLogic struct {
	ctx    context.Context
	svcCtx *rpcsvc.ServiceContext
	logx.Logger
}

func NewOnlinePushMsgLogic(ctx context.Context, svcCtx *rpcsvc.ServiceContext) *OnlinePushMsgLogic {
	return &OnlinePushMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OnlinePushMsgLogic) OnlinePushMsg(message *pb.OnlinePushMsgReq) (*pb.OnlinePushMsgResp, error) {
	var err error
	if message.MsgData.ConversationType == types.GroupChatType {
		err = l.sendMessageToGroupPush(message, message.MsgData.ReceiveID)
	} else {
		err = l.sendMessageToPush(message, message.PushToUserID)
	}
	if err != nil {
		return &pb.OnlinePushMsgResp{
			Success: false,
			ErrMsg:  err.Error(),
		}, err
	} else {
		return &pb.OnlinePushMsgResp{Success: true}, nil
	}
}

func (l *OnlinePushMsgLogic) sendMessageToGroupPush(message *pb.OnlinePushMsgReq, groupId string) error {
	// 获取群里所有没屏蔽群消息的人
	req := imuserpb.GetUserListFromGroupWithOptReq{
		GroupID: groupId,
		Opts: []imuserpb.RecvMsgOpt{
			imuserpb.RecvMsgOpt_ReceiveMessage,
			imuserpb.RecvMsgOpt_ReceiveNotNotifyMessage,
		},
		UserIDList: nil,
	}
	resp, err := l.svcCtx.ImUserRpc().GetUserListFromGroupWithOpt(l.ctx, &req)
	if err != nil {
		logx.WithContext(l.ctx).Error("GetUserListFromGroupWithOpt err :", err)
		return err
	}
	var uids []string
	for _, u := range resp.UserIDOptList {
		uids = append(uids, u.UserID)
	}
	mqPushMsg := chatpb.PushMsgDataToMQ{MsgData: message.MsgData, TraceId: xtrace.TraceIdFromContext(l.ctx), PushToUserID: uids}
	pid, offset, err := l.svcCtx.GroupPushProducer.SendMessage(l.ctx, &mqPushMsg)
	if err != nil {
		logx.WithContext(l.ctx).Error("kafka send failed ", "send data ", mqPushMsg.String(), " pid ", pid, " offset ", offset, " err ", err.Error())
		return err
	}
	return nil
}

func (l *OnlinePushMsgLogic) sendMessageToPush(message *pb.OnlinePushMsgReq, pushToUserID string) error {
	mqPushMsg := chatpb.PushMsgDataToMQ{MsgData: message.MsgData, PushToUserID: []string{pushToUserID}, TraceId: xtrace.TraceIdFromContext(l.ctx)}
	pid, offset, err := l.svcCtx.SinglePushProducer.SendMessage(l.ctx, &mqPushMsg)
	if err != nil {
		logx.WithContext(l.ctx).Error("kafka send failed", mqPushMsg.TraceId, "send data", mqPushMsg.String(), "pid", pid, "offset", offset, "err", err.Error())
		return err
	}
	return nil
}
