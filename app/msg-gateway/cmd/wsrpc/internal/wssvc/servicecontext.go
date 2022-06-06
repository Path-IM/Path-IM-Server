package wssvc

import (
	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/imuserservice"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wsconfig"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/chat"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        wsconfig.Config
	imUserService imuserservice.ImUserService
	msgRpc        chat.Chat
}

func NewServiceContext(c wsconfig.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

func (s *ServiceContext) ImUserService() imuserservice.ImUserService {
	if s.imUserService == nil {
		s.imUserService = imuserservice.NewImUserService(zrpc.MustNewClient(s.Config.ImUserRpc))
	}
	return s.imUserService
}

func (s *ServiceContext) MsgRpc() chat.Chat {
	if s.msgRpc == nil {
		s.msgRpc = chat.NewChat(zrpc.MustNewClient(s.Config.MsgRpc))
	}
	return s.msgRpc
}
