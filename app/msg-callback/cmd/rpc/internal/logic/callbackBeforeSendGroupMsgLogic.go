package logic

import (
	"context"
	"errors"

	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackBeforeSendGroupMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallbackBeforeSendGroupMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackBeforeSendGroupMsgLogic {
	return &CallbackBeforeSendGroupMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallbackBeforeSendGroupMsgLogic) CallbackBeforeSendGroupMsg(in *pb.CallbackSendGroupMsgReq) (*pb.CommonCallbackResp, error) {
	return &pb.CommonCallbackResp{}, errors.New("已弃用")
}
