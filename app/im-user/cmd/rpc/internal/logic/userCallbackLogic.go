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
	switch in.Event {
	case pb.UserCallbackReq_Online:
		l.Infof("UserCallback: online, userID: %s, platform: %s, userIp: %s", in.UserID, in.Platform, in.RemoteAddr)
	case pb.UserCallbackReq_Offline:
		l.Infof("UserCallback: offline, userID: %s, platform: %s, userIp: %s", in.UserID, in.Platform, in.RemoteAddr)
	case pb.UserCallbackReq_MsgTooOften:
		l.Infof("UserCallback: msg too often, userID: %s, platform: %s, userIp: %s", in.UserID, in.Platform, in.RemoteAddr)
	}
	return &pb.UserCallbackResp{}, nil
}
