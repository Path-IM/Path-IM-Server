package logic

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/internal/repository"

	"github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/internal/svc"
	"github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSuperGroupMaxAndMinSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	rep *repository.Rep
}

func NewGetSuperGroupMaxAndMinSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSuperGroupMaxAndMinSeqLogic {
	return &GetSuperGroupMaxAndMinSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		rep:    repository.NewRep(svcCtx),
	}
}

func (l *GetSuperGroupMaxAndMinSeqLogic) GetSuperGroupMaxAndMinSeq(in *pb.GetMaxAndMinSuperGroupSeqReq) (*pb.GetMaxAndMinSuperGroupSeqResp, error) {
	var seqs []*pb.GetMaxAndMinSuperGroupSeqRespItem
	for _, groupID := range in.SuperGroupIDList {
		maxSeq, err := l.rep.GetSuperGroupSeq(groupID)
		if err != nil {
			if err == redis.Nil {
				err = nil
			} else {
				return nil, err
			}
		}
		seqs = append(seqs, &pb.GetMaxAndMinSuperGroupSeqRespItem{
			SuperGroupID: groupID,
			MaxSeq:       uint32(maxSeq),
			MinSeq:       0,
		})
	}
	resp := new(pb.GetMaxAndMinSuperGroupSeqResp)
	resp.SuperGroupSeqList = seqs
	return resp, nil
}
