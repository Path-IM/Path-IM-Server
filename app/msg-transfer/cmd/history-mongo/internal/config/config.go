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
	Kafka     KafkaConfig
	Redis     RedisConfig
	Mongo     MongoConfig
	ImUserRpc zrpc.RpcClientConf
}
type StorageConsumer struct {
	xkafka.ProducerConfig
	MsgToHistoryGroupID string
}
type KafkaConfig struct {
	StorageConsumer StorageConsumer
	SinglePush      xkafka.ProducerConfig
	GroupPush       xkafka.ProducerConfig
}
type RedisConfig struct {
	Conf redis.RedisConf
	DB   int
}
type MongoConfig struct {
	global.MongoConfig
	DBDatabase                  string
	DBTimeout                   int
	SingleChatMsgCollectionName string
	GroupChatMsgCollectionName  string
}
