package rpclogic

import (
	"context"
	"github.com/showurl/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wslogic"

	"github.com/showurl/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcsvc"
	"github.com/showurl/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickUserConnsLogic struct {
	ctx    context.Context
	svcCtx *rpcsvc.ServiceContext
	logx.Logger
}

func NewKickUserConnsLogic(ctx context.Context, svcCtx *rpcsvc.ServiceContext) *KickUserConnsLogic {
	return &KickUserConnsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *KickUserConnsLogic) KickUserConns(in *pb.KickUserConnsReq) (*pb.KickUserConnsResp, error) {
	logic := wslogic.NewMsggatewayLogic(nil, nil)
	for _, platform := range in.PlatformIDs {
		logic.DelUserConn(in.UserID, platform)
	}
	return &pb.KickUserConnsResp{}, nil
}
