package config

import (
	"github.com/Path-IM/Path-IM-Server/common/xcql"
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf
	Kafka     KafkaConfig
	Redis     RedisConfig
	Cassandra CassandraConfig
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
type CassandraConfig struct {
	xcql.CassandraConfig
	SingleChatMsgTableName string
	GroupChatMsgTableName  string
}
