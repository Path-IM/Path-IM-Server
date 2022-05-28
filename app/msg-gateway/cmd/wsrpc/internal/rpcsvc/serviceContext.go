package rpcsvc

import (
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcconfig"
)

type ServiceContext struct {
	Config rpcconfig.Config
}

func NewServiceContext(c rpcconfig.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
