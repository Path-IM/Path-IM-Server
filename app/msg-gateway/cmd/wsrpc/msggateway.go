package main

import (
	"flag"
	"fmt"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/handler"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcconfig"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcserver"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcsvc"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wsconfig"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wssvc"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)
import _ "net/http/pprof"

var wsConfigFile = flag.String("w", "etc/msggateway-ws.yaml", "ws config file")
var rpcConfigFile = flag.String("r", "etc/msggateway-rpc.yaml", "rpc config file")

func ws() {
	flag.Parse()

	var wsConfig wsconfig.Config
	conf.MustLoad(*wsConfigFile, &wsConfig)

	ctx := wssvc.NewServiceContext(wsConfig)
	server := rest.MustNewServer(wsConfig.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", wsConfig.Host, wsConfig.Port)

	server.Start()
}

func rpc() {
	flag.Parse()

	var c rpcconfig.Config
	conf.MustLoad(*rpcConfigFile, &c)
	ctx := rpcsvc.NewServiceContext(c)
	svr := rpcserver.NewOnlineMessageRelayServiceServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterOnlineMessageRelayServiceServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:13170", nil)
	}()
	go ws()
	logx.Info("ws 启动成功 等待1秒启动 rpc")
	time.Sleep(time.Second)
	rpc()
}
