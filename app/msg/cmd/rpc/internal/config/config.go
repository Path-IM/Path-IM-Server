package config

import (
	"github.com/Path-IM/Path-IM-Server/common/xcql"
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/Path-IM/Path-IM-Server/common/xmgo/global"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Kafka          xkafka.ProducerConfig
	Callback       CallbackConfig
	MessageVerify  MessageVerifyConfig
	ImUserRpc      zrpc.RpcClientConf
	MsgCallbackRpc zrpc.RpcClientConf
	RedisConfig    RedisConfig
	Mongo          MongoConfig
	Cassandra      CassandraConfig
	HistoryDBType  string // mongo or cassandra
}

type RedisConfig struct {
	Conf redis.RedisConf
	DB   int
}
type CallbackConfig struct {
	CallbackBeforeSendGroupMsg  CallbackConfigItem
	CallbackAfterSendGroupMsg   CallbackConfigItem
	CallbackBeforeSendSingleMsg CallbackConfigItem
	CallbackAfterSendSingleMsg  CallbackConfigItem
}
type CallbackConfigItem struct {
	Enable          bool
	ContinueOnError bool
}
type MessageVerifyConfig struct {
	FriendVerify bool // 只有好友才能发送消息
}

type MongoConfig struct {
	global.MongoConfig
	DBDatabase                  string
	DBTimeout                   int
	SingleChatMsgCollectionName string
	GroupChatMsgCollectionName  string
}
type CassandraConfig struct {
	xcql.CassandraConfig
	SingleChatMsgTableName string
	GroupChatMsgTableName  string
}
