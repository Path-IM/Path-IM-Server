package rpcserver

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpclogic"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcsvc"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/attribute"
)

type kickConnConsumer struct {
	svcCtx *rpcsvc.ServiceContext
}

func (s *kickConnConsumer) HandleMsg(value []byte, key []byte, topic string, partition int32, offset int64, msg *sarama.ConsumerMessage) error {
	logx.Infof("OnlineMessageRelayServiceServer.HandleMsg: %s, %s, %d, %d", topic, string(key), partition, offset)
	msgFromMQ := &pb.KickUserConnsToMQ{}
	if err := proto.Unmarshal(msg.Value, msgFromMQ); err != nil {
		logx.Errorf("unmarshal msg error: %v", err)
		return err
	}
	var err error
	xtrace.RunWithTrace(msgFromMQ.TraceID, func(ctx context.Context) {
		logx.Infof("request kick user conns: %v", msgFromMQ.String())
		xtrace.StartFuncSpan(ctx, "MsgGateway.ConsumeKickConn.KickUserConnFromMQ", func(ctx context.Context) {
			err = s.KickUserConn(ctx, msgFromMQ)
			if err != nil {
				logx.Errorf("kick user conn error: %v", err)
			}
		})
	}, attribute.String("msg.key", string(msg.Key)))
	return err
}

func (s *kickConnConsumer) KickUserConn(ctx context.Context, msgFromMQ *pb.KickUserConnsToMQ) error {
	l := rpclogic.NewKickUserConnsLogic(ctx, s.svcCtx)
	return l.KickUserConnFromMQ(msgFromMQ)
}
