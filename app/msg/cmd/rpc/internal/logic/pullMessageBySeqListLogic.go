package logic

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/internal/repository"
	"github.com/Path-IM/Path-IM-Server/common/types"

	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullMessageBySeqListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	rep *repository.Rep
}

func NewPullMessageBySeqListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullMessageBySeqListLogic {
	return &PullMessageBySeqListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		rep:    repository.NewRep(svcCtx),
	}
}

func (l *PullMessageBySeqListLogic) PullMessageBySeqList(in *pb.PullMsgBySeqListReq) (*pb.PullMsgListResp, error) {
	resp := new(pb.PullMsgListResp)
	msgList, err := l.rep.GetMsgBySeqList(in.UserID, in.SeqList)
	if err != nil {
		l.Error("PullMessageBySeqList data error ", err.Error())
		resp.ErrCode = types.ErrCodeFailed
		resp.ErrMsg = err.Error()
		return resp, nil
	}
	resp.ErrCode = types.ErrCodeOK
	resp.ErrMsg = ""
	resp.List = msgList
	return resp, nil
}
