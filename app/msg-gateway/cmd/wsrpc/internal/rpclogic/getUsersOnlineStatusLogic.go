package rpclogic

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcsvc"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wslogic"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"

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

func (l *GetUsersOnlineStatusLogic) GetUsersOnlineStatus(req *pb.GetUsersOnlineStatusReq) (*pb.GetUsersOnlineStatusResp, error) {
	return wslogic.NewMsggatewayLogic(nil, nil).GetUsersOnlineStatus(l.ctx, req)
}
