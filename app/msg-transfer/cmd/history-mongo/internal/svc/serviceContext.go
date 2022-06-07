package svc

import (
	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/imuserservice"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-mongo/internal/config"
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	SinglePushProducer *xkafka.Producer
	GroupPushProducer  *xkafka.Producer
	imUserRpc          imuserservice.ImUserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		SinglePushProducer: xkafka.MustNewProducer(c.Kafka.SinglePush),
		GroupPushProducer:  xkafka.MustNewProducer(c.Kafka.GroupPush),
	}
}

func (s *ServiceContext) ImUserRpc() imuserservice.ImUserService {
	if s.imUserRpc == nil {
		s.imUserRpc = imuserservice.NewImUserService(zrpc.MustNewClient(s.Config.ImUserRpc))
	}
	return s.imUserRpc
}
