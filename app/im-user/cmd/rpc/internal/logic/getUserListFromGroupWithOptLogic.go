package logic

import (
	"context"
	"fmt"

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
	var uids []*pb.UserIDOpt
	for i := 0; i < 10000; i++ {
		uids = append(uids, &pb.UserIDOpt{
			UserID: fmt.Sprintf("%d", i+1),
			Opts:   pb.RecvMsgOpt_ReceiveMessage,
		})
	}
	uids = append(uids, &pb.UserIDOpt{
		UserID: "showurl",
		Opts:   pb.RecvMsgOpt_ReceiveMessage,
	}, &pb.UserIDOpt{
		UserID: "dizzy",
		Opts:   pb.RecvMsgOpt_ReceiveMessage,
	})
	return &pb.GetUserListFromGroupWithOptResp{
		CommonResp: &pb.CommonResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
		UserIDOptList: uids,
	}, nil
}
