package server

import (
	"context"
	"errors"
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
	"time"
)

func NewMsgTransferHistoryServer(svcCtx *svc.ServiceContext) *MsgTransferHistoryServer {
	m := &MsgTransferHistoryServer{svcCtx: svcCtx}
	m.cmdCh = make(chan Cmd2Value, 10000)
	m.w = new(sync.Mutex)
	m.msgHandle = make(map[string]fcb)
	m.msgHandle[svcCtx.Config.Kafka.Online.Topic] = m.ChatMs2Mongo
	m.historyConsumerGroup = xkafka.NewMConsumerGroup(&xkafka.MConsumerGroupConfig{
		KafkaVersion:   sarama.V0_10_2_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
	}, []string{svcCtx.Config.Kafka.Online.Topic},
		svcCtx.Config.Kafka.Online.Brokers, svcCtx.Config.Kafka.Online.MsgToMongoGroupID)
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
	logx.Info("msgFromMQ.OperationID: ", msgFromMQ.OperationID)
	xtrace.RunWithTrace(msgFromMQ.OperationID, func(ctx context.Context) {
		err = logic.NewMsgTransferHistoryOnlineLogic(ctx, s.svcCtx).ChatMs2Mongo(msg, msgKey)
	}, attribute.String("msgKey", msgKey))
	return err
}

func (s *MsgTransferHistoryServer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		s.SetOnlineTopicStatus(OnlineTopicBusy)
		err := s.msgHandle[msg.Topic](msg.Value, string(msg.Key))
		if err != nil {
			logx.Errorf("handle msg error: %v", err)
			continue
		}
		sess.MarkMessage(msg, "")
		if claim.HighWaterMarkOffset()-msg.Offset <= 1 {
			s.SetOnlineTopicStatus(OnlineTopicVacancy)
			s.TriggerCmd(context.Background(), OnlineTopicVacancy)
		}
	}
	return nil
}

func (s *MsgTransferHistoryServer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *MsgTransferHistoryServer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *MsgTransferHistoryServer) SetOnlineTopicStatus(busy int) {
	s.w.Lock()
	defer s.w.Unlock()
	s.OnlineTopicStatus = busy

}
func (s *MsgTransferHistoryServer) TriggerCmd(ctx context.Context, status int) {
	for {
		err := s.sendCmd(ctx, s.cmdCh, Cmd2Value{Cmd: status, Value: ""}, 1)
		if err != nil {
			logx.WithContext(ctx).Errorf("send cmd error: %v", err)
			continue
		}
		return
	}
}
func (s *MsgTransferHistoryServer) sendCmd(ctx context.Context, ch chan Cmd2Value, value Cmd2Value, timeout int64) error {
	var flag = 0
	select {
	case ch <- value:
		flag = 1
	case <-time.After(time.Second * time.Duration(timeout)):
		flag = 2
	}
	if flag == 1 {
		return nil
	} else {
		return errors.New("send cmd timeout")
	}
}
