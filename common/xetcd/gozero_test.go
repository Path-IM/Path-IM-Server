package xetcd

import (
	"context"
	"github.com/showurl/Zero-IM-Server/app/msg-gateway/cmd/wsrpc/onlineMessageRelayService"
	"github.com/showurl/Zero-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	"github.com/zeromicro/go-zero/core/discov"
	"testing"
)

func TestGetGoZeroZrpcConns(t *testing.T) {
	zrpcConns, err := GetGoZeroZrpcConns(context.TODO(), NewClient(discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
	}), "msggateway-rpc")
	if err != nil {
		t.Error(err)
	}
	for _, conn := range zrpcConns {
		onlineStatus, err := onlinemessagerelayservice.NewOnlineMessageRelayService(conn).GetUsersOnlineStatus(context.TODO(), &pb.GetUsersOnlineStatusReq{
			UserIDList:  []string{"1", "2"},
			OperationID: "",
			OpUserID:    "",
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(onlineStatus.ErrCode)
	}
}
