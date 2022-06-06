package logic

import (
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"
)

func (l *SendMsgLogic) sendMsgToKafka(m *chatpb.MsgDataToMQ, key string) error {
	m.TraceId = xtrace.TraceIdFromContext(l.ctx)
	pid, offset, err := l.svcCtx.Producer.SendMessage(l.ctx, m, key)
	if err != nil {
		l.Logger.Error(m.TraceId, " kafka send failed ", "send data ", m.String(), "pid ", pid, "offset ", offset, "err ", err.Error(), "key ", key)
	}
	return err
}
