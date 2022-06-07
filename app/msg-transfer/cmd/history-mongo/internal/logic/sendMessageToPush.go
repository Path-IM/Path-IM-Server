package logic

import (
	"context"
	imuserpb "github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
)

func (l *MsgTransferHistoryLogic) sendMessageToGroupPush(ctx context.Context, message *chatpb.MsgDataToMQ, groupId string) {
	// 获取群里所有没屏蔽群消息的人
	req := imuserpb.GetUserListFromGroupWithOptReq{
		GroupID: groupId,
		Opts: []imuserpb.RecvMsgOpt{
			imuserpb.RecvMsgOpt_ReceiveMessage,
			imuserpb.RecvMsgOpt_ReceiveNotNotifyMessage,
		},
		UserIDList: nil,
	}
	resp, err := l.svcCtx.ImUserRpc().GetUserListFromGroupWithOpt(ctx, &req)
	if err != nil {
		logx.WithContext(ctx).Error("GetUserListFromGroupWithOpt err :", err)
		return
	}
	var uids []string
	for _, u := range resp.UserIDOptList {
		uids = append(uids, u.UserID)
	}
	mqPushMsg := chatpb.PushMsgDataToMQ{MsgData: message.MsgData, TraceId: xtrace.TraceIdFromContext(l.ctx), PushToUserID: uids}
	pid, offset, err := l.svcCtx.GroupPushProducer.SendMessage(ctx, &mqPushMsg)
	if err != nil {
		logx.WithContext(ctx).Error("kafka send failed ", "send data ", mqPushMsg.String(), " pid ", pid, " offset ", offset, " err ", err.Error())
	}
}

func (l *MsgTransferHistoryLogic) sendMessageToPush(ctx context.Context, message *chatpb.MsgDataToMQ, pushToUserID string) {
	logx.WithContext(ctx).Info("msg_transfer send message to push", "message", message.String())
	mqPushMsg := chatpb.PushMsgDataToMQ{MsgData: message.MsgData, PushToUserID: []string{pushToUserID}, TraceId: xtrace.TraceIdFromContext(l.ctx)}
	pid, offset, err := l.svcCtx.SinglePushProducer.SendMessage(ctx, &mqPushMsg)
	if err != nil {
		logx.WithContext(ctx).Error("kafka send failed", mqPushMsg.TraceId, "send data", mqPushMsg.String(), "pid", pid, "offset", offset, "err", err.Error())
	}
}
