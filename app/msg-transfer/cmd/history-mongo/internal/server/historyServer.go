package server

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-mongo/internal/logic"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-mongo/internal/svc"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/attribute"
	"sync"
)

func NewMsgTransferHistoryServer(svcCtx *svc.ServiceContext) *MsgTransferHistoryServer {
	m := &MsgTransferHistoryServer{svcCtx: svcCtx}
	m.cmdCh = make(chan Cmd2Value, 10000)
	m.w = new(sync.Mutex)
	m.msgHandle = make(map[string]fcb)
	m.msgHandle[svcCtx.Config.Kafka.StorageConsumer.Topic] = m.ChatMs2Mongo
	m.historyConsumerGroup = xkafka.NewMConsumerGroup(&xkafka.MConsumerGroupConfig{
		KafkaVersion:   sarama.V0_10_2_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: true,
	}, []string{svcCtx.Config.Kafka.StorageConsumer.Topic},
		svcCtx.Config.Kafka.StorageConsumer.Brokers, svcCtx.Config.Kafka.StorageConsumer.MsgToHistoryGroupID)
	return m
}

func (s *MsgTransferHistoryServer) Start() {
	s.historyConsumerGroup.RegisterHandleAndConsumer(s)
}

func (s *MsgTransferHistoryServer) ChatMs2Mongo(msg []byte, msgKey string) error {
	msgFromMQ := chatpb.MsgDataToMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		logx.Errorf("unmarshal msg failed, err: %v", err)
		return nil
	}
	logx.Info("msgFromMQ.TraceId: ", msgFromMQ.TraceId)
	xtrace.RunWithTrace(msgFromMQ.TraceId, func(ctx context.Context) {
		err = logic.NewMsgTransferHistoryLogic(ctx, s.svcCtx).ChatMs2Mongo(msg, msgKey)
	}, attribute.String("msgKey", msgKey))
	return err
}

func (s *MsgTransferHistoryServer) HandleMsg(value []byte, key []byte, topic string, partition int32, offset int64, msg *sarama.ConsumerMessage) error {
	err := s.msgHandle[msg.Topic](msg.Value, string(msg.Key))
	if err != nil {
		logx.Errorf("handle msg error: %v", err)
		return err
	}
	return nil
}
