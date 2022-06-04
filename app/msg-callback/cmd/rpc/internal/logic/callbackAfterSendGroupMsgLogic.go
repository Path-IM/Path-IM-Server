package logic

import (
	"context"
	"errors"

	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/pb"

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
	return &pb.CommonCallbackResp{}, errors.New("已弃用")
}
