package main

import (
	"flag"
	"github.com/showurl/Path-IM-Server/app/msg-transfer/cmd/history/internal/config"
	"github.com/showurl/Path-IM-Server/app/msg-transfer/cmd/history/internal/server"
	"github.com/showurl/Path-IM-Server/app/msg-transfer/cmd/history/internal/svc"
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
