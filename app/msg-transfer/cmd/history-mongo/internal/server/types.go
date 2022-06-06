package server

import (
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-mongo/internal/svc"
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"sync"
)

type fcb func(msg []byte, msgKey string) error
type Cmd2Value struct {
	Cmd   int
	Value interface{}
}

type MsgTransferHistoryServer struct {
	svcCtx               *svc.ServiceContext
	msgHandle            map[string]fcb
	historyConsumerGroup *xkafka.MConsumerGroup
	cmdCh                chan Cmd2Value
	w                    *sync.Mutex
}
