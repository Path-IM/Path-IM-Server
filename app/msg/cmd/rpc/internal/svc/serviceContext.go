package svc

import (
	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/imuserservice"
	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/msgcallbackservice"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/internal/config"
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	ImUser      imuserservice.ImUserService
	MsgCallback msgcallbackservice.MsgcallbackService

	Producer *xkafka.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		ImUser:      imuserservice.NewImUserService(zrpc.MustNewClient(c.ImUserRpc)),
		MsgCallback: msgcallbackservice.NewMsgcallbackService(zrpc.MustNewClient(c.MsgCallbackRpc)),
		Producer:    xkafka.MustNewProducer(c.Kafka),
	}
}
