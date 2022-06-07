package rpcsvc

import (
	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/imuserservice"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcconfig"
	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/msgpushservice"
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             rpcconfig.Config
	imUserRpc          imuserservice.ImUserService
	SinglePushProducer *xkafka.Producer
	GroupPushProducer  *xkafka.Producer
	msgPushRpc         msgpushservice.MsgPushService
}

func (c *ServiceContext) ImUserRpc() imuserservice.ImUserService {
	if c.imUserRpc == nil {
		c.imUserRpc = imuserservice.NewImUserService(zrpc.MustNewClient(c.Config.ImUserRpc))
	}
	return c.imUserRpc
}

func (c *ServiceContext) MsgPushRpc() msgpushservice.MsgPushService {
	if c.msgPushRpc == nil {
		c.msgPushRpc = msgpushservice.NewMsgPushService(zrpc.MustNewClient(c.Config.MsgPushRpc))
	}
	return c.msgPushRpc
}

func NewServiceContext(c rpcconfig.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		SinglePushProducer: xkafka.MustNewProducer(c.Producer.SinglePush),
		GroupPushProducer:  xkafka.MustNewProducer(c.Producer.GroupPush),
	}
}
