package wssvc

import (
	"github.com/showurl/Path-IM-Server/app/im-user/cmd/rpc/imuserservice"
	"github.com/showurl/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wsconfig"
	"github.com/showurl/Path-IM-Server/app/msg/cmd/rpc/chat"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        wsconfig.Config
	ImUserService imuserservice.ImUserService
	MsgRpc        chat.Chat
}

func NewServiceContext(c wsconfig.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		ImUserService: imuserservice.NewImUserService(zrpc.MustNewClient(c.ImUserRpc)),
		MsgRpc:        chat.NewChat(zrpc.MustNewClient(c.MsgRpc)),
	}
}
