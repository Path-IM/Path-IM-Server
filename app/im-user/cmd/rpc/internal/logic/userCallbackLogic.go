package logic

import (
	"context"

	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCallbackLogic {
	return &UserCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  用户回调
func (l *UserCallbackLogic) UserCallback(in *pb.UserCallbackReq) (*pb.UserCallbackResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UserCallbackResp{}, nil
}
