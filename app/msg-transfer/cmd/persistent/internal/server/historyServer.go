package server

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/persistent/internal/logic"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/persistent/internal/svc"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/attribute"
	"sync"
)

func NewMsgTransferPersistentServer(svcCtx *svc.ServiceContext) *MsgTransferPersistentServer {
	m := &MsgTransferPersistentServer{svcCtx: svcCtx}
	m.cmdCh = make(chan Cmd2Value, 10000)
	m.w = new(sync.Mutex)
	m.msgHandle = make(map[string]fcb)
	m.msgHandle[svcCtx.Config.Kafka.Topic] = m.SaveMsg
	m.persistentConsumerGroup = xkafka.NewMConsumerGroup(&xkafka.MConsumerGroupConfig{
		KafkaVersion:   sarama.V0_10_2_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
	}, []string{svcCtx.Config.Kafka.Topic},
		svcCtx.Config.Kafka.Brokers, svcCtx.Config.Kafka.MsgPersistentGroupID)
	return m
}

func (s *MsgTransferPersistentServer) Start() {
	s.persistentConsumerGroup.RegisterHandleAndConsumer(s)
}

func (s *MsgTransferPersistentServer) SaveMsg(msg []byte, msgKey string) error {
	msgFromMQ := chatpb.MsgDataToMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		logx.Errorf("unmarshal msg failed, err: %v", err)
		return nil
	}
	logx.Info("msgFromMQ.TraceId: ", msgFromMQ.TraceId)
	xtrace.RunWithTrace(msgFromMQ.TraceId, func(ctx context.Context) {
		err = logic.NewMsgTransferPersistentLogic(ctx, s.svcCtx).Do(msg, msgKey)
	}, attribute.String("msgKey", msgKey))
	return err
}

func (s *MsgTransferPersistentServer) HandleMsg(value []byte, key []byte, topic string, partition int32, offset int64, msg *sarama.ConsumerMessage) error {
	err := s.msgHandle[msg.Topic](msg.Value, string(msg.Key))
	if err != nil {
		logx.Errorf("msgHandle error: %v", err)
		return err
	}
	return nil
}
