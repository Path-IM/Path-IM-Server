package xetcd

import (
	"context"
	"github.com/zeromicro/go-zero/zrpc"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func GetGoZeroZrpcConns(ctx context.Context, etcdClient *clientv3.Client, key string) (zrpcConns []zrpc.Client, err error) {
	response, err := etcdClient.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	for _, kv := range response.Kvs {
		zrpcConns = append(zrpcConns, zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: []string{string(kv.Value)},
			Target:    "",
			App:       "",
			Token:     "",
			NonBlock:  false,
			Timeout:   0,
		}))
	}
	return
}
