package logic

import (
	"context"

	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListFromGroupWithOptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListFromGroupWithOptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListFromGroupWithOptLogic {
	return &GetUserListFromGroupWithOptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取群成员列表 通过消息接收选项
func (l *GetUserListFromGroupWithOptLogic) GetUserListFromGroupWithOpt(in *pb.GetUserListFromGroupWithOptReq) (*pb.GetUserListFromGroupWithOptResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserListFromGroupWithOptResp{
		CommonResp: &pb.CommonResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
		UserIDOptList: []*pb.UserIDOpt{{
			UserID: "1",
			Opts:   pb.RecvMsgOpt_ReceiveMessage,
		}},
	}, nil
}
