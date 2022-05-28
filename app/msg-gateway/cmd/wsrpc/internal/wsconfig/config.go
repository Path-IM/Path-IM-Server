package wsconfig

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	ImUserRpc        zrpc.RpcClientConf
	MsgRpc           zrpc.RpcClientConf
	Websocket        WebsocketConfig
	SendMsgRateLimit RateLimitConfig
	Redis            redis.RedisConf
}
type WebsocketConfig struct {
	MaxConnNum int
	TimeOut    int
	MaxMsgLen  int
}
type RateLimitConfig struct {
	Enable  bool
	Seconds int
	Quota   int
}
