package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/internal/config"
	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/internal/server"
	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)
import _ "net/http/pprof"

var configFile = flag.String("f", "etc/msgpush.yaml", "the config file")

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:13171", nil)
	}()
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	svr := server.NewMsgPushServiceServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterMsgPushServiceServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
