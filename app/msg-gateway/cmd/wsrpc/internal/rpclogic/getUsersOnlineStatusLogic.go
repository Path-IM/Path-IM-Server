package rpclogic

import (
	"context"
	"github.com/showurl/Zero-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcsvc"

	"github.com/showurl/Zero-IM-Server/app/msg-gateway/cmd/wsrpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersOnlineStatusLogic struct {
	ctx    context.Context
	svcCtx *rpcsvc.ServiceContext
	logx.Logger
}

func NewGetUsersOnlineStatusLogic(ctx context.Context, svcCtx *rpcsvc.ServiceContext) *GetUsersOnlineStatusLogic {
	return &GetUsersOnlineStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUsersOnlineStatusLogic) GetUsersOnlineStatus(in *pb.GetUsersOnlineStatusReq) (*pb.GetUsersOnlineStatusResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUsersOnlineStatusResp{
		ErrCode:       0,
		ErrMsg:        "",
		SuccessResult: nil,
		FailedResult:  nil,
	}, nil
}
