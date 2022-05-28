package logic

import (
	"context"

	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type IfAInBFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIfAInBFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IfAInBFriendListLogic {
	return &IfAInBFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  判断用户A是否在B好友列表中
func (l *IfAInBFriendListLogic) IfAInBFriendList(in *pb.IfAInBFriendListReq) (*pb.IfAInBFriendListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.IfAInBFriendListResp{
		CommonResp: &pb.CommonResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
		IsInFriendList: true,
	}, nil
}
