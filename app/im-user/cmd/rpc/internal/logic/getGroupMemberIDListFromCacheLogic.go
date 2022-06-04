package logic

import (
	"context"
	"errors"

	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberIDListFromCacheLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberIDListFromCacheLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberIDListFromCacheLogic {
	return &GetGroupMemberIDListFromCacheLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMemberIDListFromCacheLogic) GetGroupMemberIDListFromCache(in *pb.GetGroupMemberIDListFromCacheReq) (*pb.GetGroupMemberIDListFromCacheResp, error) {
	return &pb.GetGroupMemberIDListFromCacheResp{}, errors.New("已弃用")
}
