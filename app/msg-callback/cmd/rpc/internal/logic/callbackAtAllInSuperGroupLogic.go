package logic

import (
	"context"

	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackAtAllInSuperGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallbackAtAllInSuperGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackAtAllInSuperGroupLogic {
	return &CallbackAtAllInSuperGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallbackAtAllInSuperGroupLogic) CallbackAtAllInSuperGroup(in *pb.CallbackAtAllInSuperGroupReq) (*pb.CallbackAtAllInSuperGroupResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CallbackAtAllInSuperGroupResp{CommonCallbackResp: &pb.CommonCallbackResp{
		ActionCode: pb.ActionCode_Forbidden,
		ErrCode:    pb.ErrCode_HandleFailed,
		ErrMsg:     "",
	}}, nil
}
