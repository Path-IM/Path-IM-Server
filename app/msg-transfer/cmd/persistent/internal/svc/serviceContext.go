package svc

import (
	"github.com/showurl/Path-IM-Server/app/msg-transfer/cmd/persistent/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
