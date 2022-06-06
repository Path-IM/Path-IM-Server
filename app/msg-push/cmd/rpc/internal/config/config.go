package config

import (
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	PushType                string `json:",default=jpns,options=jpns|mobpush"`
	Jpns                    JpnsConf
	MsgGatewayRpcEtcd       *discov.EtcdConf `json:",optional"`
	MsgGatewayRpcEndpoints  []string         `json:",optional"`
	ImUserRpc               zrpc.RpcClientConf
	SinglePushConsumer      SinglePushConsumerConfig
	GroupPushConsumer       GroupPushConsumerConfig
	MsgGatewayRpcK8sTarget  string `json:",optional"`
	OfflinePushDefaultTitle string // 默认的离线推送标题
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

type GroupPushConsumerConfig struct {
	xkafka.ProducerConfig
	GroupPushGroupID string
}
