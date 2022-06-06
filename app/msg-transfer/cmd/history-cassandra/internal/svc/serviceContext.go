package svc

import (
	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/msgpushservice"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-cassandra/internal/config"
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	SinglePushProducer *xkafka.Producer
	GroupPushProducer  *xkafka.Producer
	MsgPush            msgpushservice.MsgPushService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		SinglePushProducer: xkafka.MustNewProducer(c.Kafka.SinglePush),
		GroupPushProducer:  xkafka.MustNewProducer(c.Kafka.GroupPush),
		MsgPush:            msgpushservice.NewMsgPushService(zrpc.MustNewClient(c.MsgPushRpc)),
	}
}
