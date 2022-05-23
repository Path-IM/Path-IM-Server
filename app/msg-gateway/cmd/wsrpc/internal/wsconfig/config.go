package wsconfig

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	ImUserRpc       zrpc.RpcClientConf
	MsgRpc          zrpc.RpcClientConf
	Websocket       WebsocketConfig
	VerifyTokenOnce bool
}
type WebsocketConfig struct {
	MaxConnNum int
	TimeOut    int
	MaxMsgLen  int
}
