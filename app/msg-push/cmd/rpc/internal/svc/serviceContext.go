package svc

import (
	"github.com/showurl/Zero-IM-Server/app/im-user/cmd/rpc/imuserservice"
	"github.com/showurl/Zero-IM-Server/app/msg-push/cmd/rpc/internal/config"
	"github.com/showurl/Zero-IM-Server/app/msg-push/cmd/rpc/internal/sdk"
	push "github.com/showurl/Zero-IM-Server/app/msg-push/cmd/rpc/internal/sdk/jpush"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	offlinePusher sdk.OfflinePusher
	ImUserService imuserservice.ImUserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		ImUserService: imuserservice.NewImUserService(zrpc.MustNewClient(c.ImUserRpc)),
	}
}
func (s *ServiceContext) GetOfflinePusher() sdk.OfflinePusher {
	if s.offlinePusher != nil {
		return s.offlinePusher
	}
	if s.Config.PushType == "jpns" {
		s.offlinePusher = &push.JPush{s.Config}
	}
	return s.offlinePusher
}
