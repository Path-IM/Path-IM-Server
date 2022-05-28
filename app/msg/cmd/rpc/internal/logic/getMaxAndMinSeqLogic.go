package logic

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/internal/repository"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/go-redis/redis/v8"

	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"

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
	resp := new(pb.GetMaxAndMinSeqResp)
	maxSeq, err := l.rep.GetUserMaxSeq(in.UserID)
	if err == nil {
		resp.MaxSeq = uint32(maxSeq)
	} else if err == redis.Nil {
		resp.MaxSeq = 0
	} else {
		l.Error("getMaxSeq from redis error", in.String(), err.Error())
		resp.ErrCode = types.ErrCodeFailed
		resp.ErrMsg = "redis get err"
		return resp, nil
	}
	minSeq, err := l.rep.GetUserMinSeq(in.UserID)
	if err == nil {
		resp.MinSeq = uint32(minSeq)
	} else if err == redis.Nil {
		resp.MinSeq = 0
	} else {
		l.Error("getMinSeq from redis error", in.String(), err.Error())
		resp.ErrCode = types.ErrCodeFailed
		resp.ErrMsg = "redis get err"
		return resp, nil
	}
	return resp, nil
}
