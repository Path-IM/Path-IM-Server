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

type GetMinAndMaxGroupSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	rep *repository.Rep
}

func NewGetMinAndMaxGroupSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMinAndMaxGroupSeqLogic {
	return &GetMinAndMaxGroupSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		rep:    repository.NewRep(svcCtx),
	}
}

func (l *GetMinAndMaxGroupSeqLogic) GetMinAndMaxGroupSeq(in *pb.GetMinAndMaxGroupSeqReq) (*pb.GetMinAndMaxGroupSeqResp, error) {
	resp := new(pb.GetMinAndMaxGroupSeqResp)
	var seqs []*pb.GetMinAndMaxGroupSeqItem
	for _, groupID := range in.GroupIDList {
		maxSeq, err := l.rep.GetGroupMaxSeq(groupID)
		if err != nil {
			if err == redis.Nil {
				err = nil
			} else {
				l.Error("GetGroupMaxSeq err ", err)
				resp.ErrCode = types.ErrCodeFailed
				resp.ErrMsg = err.Error()
				//return nil, err
			}
		}
		minSeq, err := l.rep.GetGroupMinSeq(groupID)
		if err != nil {
			if err == redis.Nil {
				err = nil
			} else {
				l.Error("GetGroupMaxSeq err ", err)
				resp.ErrCode = types.ErrCodeFailed
				resp.ErrMsg = err.Error()
				//return nil, err
			}
		}
		seqs = append(seqs, &pb.GetMinAndMaxGroupSeqItem{
			GroupID: groupID,
			MaxSeq:  uint32(maxSeq),
			MinSeq:  uint32(minSeq),
		})
	}
	resp.GroupSeqList = seqs
	return resp, nil
}
