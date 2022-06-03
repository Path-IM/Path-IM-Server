package config

import (
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	PushType string `json:",default=jpns,options=jpns|mobpush"`
	//MsgGatewayEtcd discov.EtcdConf
	Jpns                   JpnsConf
	MsgGatewayRpc          discov.EtcdConf
	ImUserRpc              zrpc.RpcClientConf
	SinglePushConsumer     SinglePushConsumerConfig
	SuperGroupPushConsumer SuperGroupPushConsumerConfig
	MsgGatewayRpcK8sTarget string `json:",optional"`
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
