package onlinemessagerelayservice

import (
	"context"
	"github.com/zeromicro/go-zero/core/discov"
)

// GetAllByEtcd 获取所有 service
// @param msggatewayRpcEtcdKey msggateway-rpc 的 etcd key
func GetAllByEtcd(
	ctx context.Context,
	etcdConf discov.EtcdConf,
	msggatewayRpcEtcdKey string,
) (services []OnlineMessageRelayService, err error) {
	// todo get all podip from etcd
	panic("未实现")
	return
}
