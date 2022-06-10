package rpclogic

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcsvc"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wslogic"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"

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
	//logic := wslogic.NewMsggatewayLogic(nil, nil)
	//for _, platform := range in.PlatformIDs {
	//	logic.DelUserConn(in.UserID, platform)
	//}
	msg := &pb.KickUserConnsToMQ{
		UserID:      in.UserID,
		PlatformIDs: in.PlatformIDs,
		TraceID:     xtrace.TraceIdFromContext(l.ctx),
	}
	_, _, err := l.svcCtx.KickConnProducer.SendMessage(l.ctx, msg)
	if err != nil {
		return &pb.KickUserConnsResp{
			ErrCode: 500,
			ErrMsg:  err.Error(),
		}, err
	}
	return &pb.KickUserConnsResp{}, nil
}

func (l *KickUserConnsLogic) KickUserConnFromMQ(in *pb.KickUserConnsToMQ) error {
	logic := wslogic.NewMsggatewayLogic(nil, nil)
	for _, platform := range in.PlatformIDs {
		logic.DelUserConn(in.UserID, platform)
	}
	return nil
}
