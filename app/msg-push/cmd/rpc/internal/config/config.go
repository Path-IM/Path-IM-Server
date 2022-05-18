package config

import (
	"github.com/showurl/Zero-IM-Server/common/xkafka"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	PushType string `json:",default=jpns,options=jpns|mobpush"`
	//MsgGatewayEtcd discov.EtcdConf
	Jpns                   JpnsConf
	MsgGatewayRpc          zrpc.RpcClientConf
	ImUserRpc              zrpc.RpcClientConf
	SinglePushConsumer     SinglePushConsumerConfig
	SuperGroupPushConsumer SuperGroupPushConsumerConfig
}
type JpnsConf struct {
	PushIntent     string
	PushUrl        string
	AppKey         string
	MasterSecret   string
	ApnsProduction bool `json:",default=false"`
}

type SinglePushConsumerConfig struct {
	xkafka.ProducerConfig
	SinglePushGroupID string
}

type SuperGroupPushConsumerConfig struct {
	xkafka.ProducerConfig
	SuperGroupPushGroupID string
}
