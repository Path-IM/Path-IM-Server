package logic

import (
	"context"

	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListFromSuperGroupWithOptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListFromSuperGroupWithOptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListFromSuperGroupWithOptLogic {
	return &GetUserListFromSuperGroupWithOptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取超级群成员列表 通过消息接收选项
func (l *GetUserListFromSuperGroupWithOptLogic) GetUserListFromSuperGroupWithOpt(in *pb.GetUserListFromSuperGroupWithOptReq) (*pb.GetUserListFromSuperGroupWithOptResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserListFromSuperGroupWithOptResp{
		CommonResp: &pb.CommonResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
		UserIDOptList: []*pb.UserIDOpt{{
			UserID: "1",
			Opts:   in.Opts[0],
		}, {
			UserID: "2",
			Opts:   in.Opts[0],
		}, {
			UserID: "3",
			Opts:   in.Opts[0],
		}, {
			UserID: "4",
			Opts:   in.Opts[0],
		}, {
			UserID: "5",
			Opts:   in.Opts[0],
		}, {
			UserID: "6",
			Opts:   in.Opts[0],
		}, {
			UserID: "7",
			Opts:   in.Opts[0],
		}, {
			UserID: "8",
			Opts:   in.Opts[0],
		}, {
			UserID: "9",
			Opts:   in.Opts[0],
		}, {
			UserID: "10",
			Opts:   in.Opts[0],
		}},
	}, nil
}
