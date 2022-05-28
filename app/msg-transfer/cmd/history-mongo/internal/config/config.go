package config

import (
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/Path-IM/Path-IM-Server/common/xmgo/global"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf
	Kafka      KafkaConfig
	Redis      RedisConfig
	Mongo      MongoConfig
	MsgPushRpc zrpc.RpcClientConf
}
type KafkaConfigOnline struct {
	xkafka.ProducerConfig
	MsgToMongoGroupID string
}
type KafkaConfig struct {
	Online         KafkaConfigOnline
	Offline        xkafka.ProducerConfig
	SinglePush     xkafka.ProducerConfig
	SuperGroupPush xkafka.ProducerConfig
}
type RedisConfig struct {
	Conf redis.RedisConf
	DB   int
}
type MongoConfig struct {
	global.MongoConfig
	DBDatabase                      string
	DBTimeout                       int
	SingleChatMsgCollectionName     string
	SuperGroupChatMsgCollectionName string
}
