package logic

import (
	"context"

	"github.com/showurl/Path-IM-Server/app/im-user/cmd/rpc/internal/svc"
	"github.com/showurl/Path-IM-Server/app/im-user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type IfAInBBlacklistLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIfAInBBlacklistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IfAInBBlacklistLogic {
	return &IfAInBBlacklistLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  判断用户A是否在B黑名单中
func (l *IfAInBBlacklistLogic) IfAInBBlacklist(in *pb.IfAInBBlacklistReq) (*pb.IfAInBBlacklistResp, error) {
	// todo: add your logic here and delete this line

	return &pb.IfAInBBlacklistResp{
		CommonResp: &pb.CommonResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
		IsInBlacklist: false,
	}, nil
}
