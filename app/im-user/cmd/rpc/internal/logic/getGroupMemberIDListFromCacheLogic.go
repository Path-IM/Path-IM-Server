package logic

import (
	"context"

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
	// todo: add your logic here and delete this line
	// 如果你使用 Open-IM 的群聊功能 此处需要你实现
	// 如果你仅仅使用 Zero-IM 的超级大群功能 你需要实现 GetUserListFromSuperGroupWithOpt rpc接口
	return &pb.GetGroupMemberIDListFromCacheResp{
		CommonResp: &pb.CommonResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
		MemberIDList: []string{
			"1",
			"2",
			"3",
			"4",
			"5",
			"6",
			"7",
			"8",
			"9",
			"10",
		},
	}, nil
}
