package config

import (
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	PushType                string      `json:",default=jpns,options=jpns|mobpush"`
	Jpns                    JpnsConf    `json:",optional"`
	MobPush                 MobPushConf `json:",optional"`
	ImUserRpc               zrpc.RpcClientConf
	OfflinePushDefaultTitle string // 默认的离线推送标题
	OfflinePushGroupTitle   string // 群聊消息推送标题
	SinglePushConsumer      SinglePushConsumerConfig
	GroupPushConsumer       GroupPushConsumerConfig
	Redis                   redis.RedisConf
}
type JpnsConf struct {
	PushIntent     string
	PushUrl        string
	AppKey         string
	MasterSecret   string
	ApnsProduction bool `json:",default=false"`
}
type MobPushConf struct {
	AppKey         string
	AppSecret      string
	ApnsProduction bool `json:",default=false"`
	ApnsCateGory   string
	ApnsSound      string
	AndroidSound   string
}

type SinglePushConsumerConfig struct {
	xkafka.ProducerConfig
	SinglePushGroupID string
}

type GroupPushConsumerConfig struct {
	xkafka.ProducerConfig
	GroupPushGroupID string
}
