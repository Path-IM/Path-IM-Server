package logic

import (
	"context"

	"github.com/showurl/Path-IM-Server/app/msg-callback/cmd/rpc/internal/svc"
	"github.com/showurl/Path-IM-Server/app/msg-callback/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackAfterSendGroupMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallbackAfterSendGroupMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackAfterSendGroupMsgLogic {
	return &CallbackAfterSendGroupMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallbackAfterSendGroupMsgLogic) CallbackAfterSendGroupMsg(in *pb.CallbackSendGroupMsgReq) (*pb.CommonCallbackResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CommonCallbackResp{}, nil
}
