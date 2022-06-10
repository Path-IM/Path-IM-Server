package rpcconfig

import (
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
)

type Config struct {
	zrpc.RpcServerConf
	SinglePushConsumer    SinglePushConsumerConfig
	GroupPushConsumer     GroupPushConsumerConfig
	KickConnConsumerGroup KickConnConsumerConfig
	ImUserRpc             zrpc.RpcClientConf
	MsgPushRpc            zrpc.RpcClientConf
	Producer              KafkaConfig
}

type KafkaConfig struct {
	SinglePush xkafka.ProducerConfig
	GroupPush  xkafka.ProducerConfig
	KickConn   xkafka.ProducerConfig
}
type SinglePushConsumerConfig struct {
	xkafka.ProducerConfig
	SinglePushGroupID string `json:",optional"`
}

type GroupPushConsumerConfig struct {
	xkafka.ProducerConfig
	GroupPushGroupID string `json:",optional"`
}
type KickConnConsumerConfig struct {
	xkafka.ProducerConfig
	KickConnGroupID string `json:",optional"`
}

func (s SinglePushConsumerConfig) GetGroupID() string {
	if s.SinglePushGroupID == "" {
		podName := os.Getenv("POD_NAME")
		if podName == "" {
			panic("env POD_NAME is null")
		}
		return podName
	}
	return s.SinglePushGroupID
}
func (s GroupPushConsumerConfig) GetGroupID() string {
	if s.GroupPushGroupID == "" {
		podName := os.Getenv("POD_NAME")
		if podName == "" {
			panic("env POD_NAME is null")
		}
		return podName
	}
	return s.GroupPushGroupID
}
func (s KickConnConsumerConfig) GetGroupID() string {
	if s.KickConnGroupID == "" {
		podName := os.Getenv("POD_NAME")
		if podName == "" {
			panic("env POD_NAME is null")
		}
		return podName
	}
	return s.KickConnGroupID
}
