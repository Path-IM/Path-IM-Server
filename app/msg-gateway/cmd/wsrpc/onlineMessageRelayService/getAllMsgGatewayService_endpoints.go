package onlinemessagerelayservice

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/zrpc"
)

var (
	servicesByEndpoints []OnlineMessageRelayService
)

// GetAllByEndpoints 获取所有 service
// @param msggatewayRpcEtcdKey msggateway-rpc 的 etcd key
func GetAllByEndpoints(
	ctx context.Context,
	endpoints []string,
) (services []OnlineMessageRelayService, err error) {
	if endpoints == nil {
		return nil, errors.New("endpoints is nil")
	}
	if servicesByEndpoints == nil {
		for _, endpoint := range endpoints {
			servicesByEndpoints = append(servicesByEndpoints, NewOnlineMessageRelayService(
				zrpc.MustNewClient(zrpc.RpcClientConf{
					Endpoints: []string{endpoint},
				}),
			))
		}
	}
	return servicesByEndpoints, nil
}
