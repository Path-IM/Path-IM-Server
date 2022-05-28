package logic

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/internal/repository"

	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullMessageBySuperGroupSeqListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	rep *repository.Rep
}

func NewPullMessageBySuperGroupSeqListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullMessageBySuperGroupSeqListLogic {
	return &PullMessageBySuperGroupSeqListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		rep:    repository.NewRep(svcCtx),
	}
}

func (l *PullMessageBySuperGroupSeqListLogic) PullMessageBySuperGroupSeqList(in *pb.PullMessageBySuperGroupSeqListReq) (*pb.WrapPullMessageBySeqListResp, error) {
	resp := new(pb.PullMessageBySeqListResp)
	msgList, err := l.rep.GetMsgBySuperGroupSeqList(in.GroupID, in.SeqList)
	if err != nil {
		l.Error("PullMessageBySeqList data error ", err.Error())
		resp.ErrCode = 201
		resp.ErrMsg = err.Error()
		return &pb.WrapPullMessageBySeqListResp{
			PullMessageBySeqListResp: resp,
		}, nil
	}
	//respSingleMsgFormat = singleMsgHandleByUser(SingleMsgFormat, in.UserID)
	//respGroupMsgFormat = groupMsgHandleByUser(GroupMsgFormat)
	resp.ErrCode = 0
	resp.ErrMsg = ""
	resp.List = msgList
	return &pb.WrapPullMessageBySeqListResp{
		PullMessageBySeqListResp: resp,
	}, nil
}
