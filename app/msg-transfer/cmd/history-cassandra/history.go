package main

import (
	"flag"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-cassandra/internal/config"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-cassandra/internal/server"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-cassandra/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/history.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	err := c.ServiceConf.SetUp()
	if err != nil {
		panic(err)
	}
	s := server.NewMsgTransferHistoryServer(svc.NewServiceContext(c))
	s.Start()
}
