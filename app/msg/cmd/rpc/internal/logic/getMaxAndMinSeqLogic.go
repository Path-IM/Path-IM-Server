package logic

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/internal/repository"

	"github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/internal/svc"
	"github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMaxAndMinSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	rep *repository.Rep
}

func NewGetMaxAndMinSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMaxAndMinSeqLogic {
	return &GetMaxAndMinSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		rep:    repository.NewRep(svcCtx),
	}
}

func (l *GetMaxAndMinSeqLogic) GetMaxAndMinSeq(in *pb.GetMaxAndMinSeqReq) (*pb.GetMaxAndMinSeqResp, error) {
	maxSeq, err := l.rep.GetUserSeq(in.UserID)
	resp := new(pb.GetMaxAndMinSeqResp)
	if err == nil {
		resp.MaxSeq = uint32(maxSeq)
	} else if err == redis.Nil {
		resp.MaxSeq = 0
	} else {
		l.Error("getMaxSeq from redis error", in.String(), err.Error())
		resp.ErrCode = 200
		resp.ErrMsg = "redis get err"
	}
	resp.MinSeq = 0
	return resp, nil
}
