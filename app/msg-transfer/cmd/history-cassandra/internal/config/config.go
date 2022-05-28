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
	Kafka      KafkaConfig
	Redis      RedisConfig
	Cassandra  CassandraConfig
	MsgPushRpc zrpc.RpcClientConf
}
type KafkaConfigOnline struct {
	xkafka.ProducerConfig
	MsgToCassandraGroupID string
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
type CassandraConfig struct {
	xcql.CassandraConfig
	DBDatabase                 string
	SingleChatMsgTableName     string
	SuperGroupChatMsgTableName string
}
