package logic

import (
	"errors"
	chatpb "github.com/showurl/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/showurl/Path-IM-Server/common/types"
	"github.com/showurl/Path-IM-Server/common/xtrace"
)

func (l *SendMsgLogic) sendMsgToKafka(m *chatpb.MsgDataToMQ, key string, status string) error {
	m.OperationID = xtrace.TraceIdFromContext(l.ctx)
	switch status {
	case types.OnlineStatus:
		pid, offset, err := l.svcCtx.OnlineProducer.SendMessage(l.ctx, m, key)
		if err != nil {
			l.Logger.Error(m.OperationID, "kafka send failed ", "send data ", m.String(), "pid ", pid, "offset ", offset, "err ", err.Error(), "key ", key, status)
		}
		return err
	case types.OfflineStatus:
		pid, offset, err := l.svcCtx.OfflineProducer.SendMessage(l.ctx, m, key)
		if err != nil {
			l.Logger.Error(m.OperationID, "kafka send failed ", "send data ", m.String(), "pid ", pid, "offset ", offset, "err ", err.Error(), "key ", key, status)
		}
		return err
	}
	return errors.New("status error")
}
