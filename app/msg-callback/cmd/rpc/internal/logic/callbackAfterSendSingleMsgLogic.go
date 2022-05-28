package logic

import (
	"context"

	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackAfterSendSingleMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallbackAfterSendSingleMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackAfterSendSingleMsgLogic {
	return &CallbackAfterSendSingleMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallbackAfterSendSingleMsgLogic) CallbackAfterSendSingleMsg(in *pb.CallbackSendSingleMsgReq) (*pb.CommonCallbackResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CommonCallbackResp{
		ActionCode: pb.ActionCode_Forbidden,
		ErrCode:    pb.ErrCode_HandleFailed,
		ErrMsg:     "",
	}, nil
}
