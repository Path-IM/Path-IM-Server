package rpcserver

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpclogic"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcsvc"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/attribute"
)

type groupConsumer struct {
	svcCtx *rpcsvc.ServiceContext
}

func (s *groupConsumer) HandleMsg(value []byte, key []byte, topic string, partition int32, offset int64, msg *sarama.ConsumerMessage) error {
	logx.Infof("OnlineMessageRelayServiceServer.HandleMsg: %s, %s, %d, %d", topic, string(key), partition, offset)
	msgFromMQ := &chatpb.PushMsgDataToMQ{}
	if err := proto.Unmarshal(msg.Value, msgFromMQ); err != nil {
		logx.Errorf("unmarshal msg error: %v", err)
		return err
	}
	var err error
	xtrace.RunWithTrace(msgFromMQ.TraceId, func(ctx context.Context) {
		xtrace.StartFuncSpan(ctx, "MsgGateway.ConsumeGroup.PushMsg2GroupMember", func(ctx context.Context) {
			err = s.PushMsg(ctx, msgFromMQ)
			if err != nil {
				logx.Errorf("push Group msg error: %v", err)
			}
		})
	}, attribute.String("msg.key", string(msg.Key)))
	return err
}

func (s *groupConsumer) PushMsg(ctx context.Context, msgFromMQ *chatpb.PushMsgDataToMQ) error {
	l := rpclogic.NewPushMsgLogic(ctx, s.svcCtx)
	return l.PushMsg(msgFromMQ)
}
