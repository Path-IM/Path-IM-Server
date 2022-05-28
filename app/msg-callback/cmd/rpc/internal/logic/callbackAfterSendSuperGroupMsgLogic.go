package logic

import (
	"context"

	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackAfterSendSuperGroupMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallbackAfterSendSuperGroupMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackAfterSendSuperGroupMsgLogic {
	return &CallbackAfterSendSuperGroupMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallbackAfterSendSuperGroupMsgLogic) CallbackAfterSendSuperGroupMsg(in *pb.CallbackSendSuperGroupMsgReq) (*pb.CommonCallbackResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CommonCallbackResp{
		ActionCode: 0,
		ErrCode:    0,
		ErrMsg:     "",
	}, nil
}
