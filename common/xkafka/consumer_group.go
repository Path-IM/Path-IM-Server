package xkafka

import (
	"context"
	"github.com/Shopify/sarama"
)

type MConsumerGroup struct {
	sarama.ConsumerGroup
	groupID string
	topics  []string
}
type msgConsumer interface {
	HandleMsg(value []byte, key []byte, topic string, partition int32, offset int64, msg *sarama.ConsumerMessage) error
}
type consumer struct {
	msgConsumer
}

func (c consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := c.HandleMsg(msg.Value, msg.Key, msg.Topic, msg.Partition, msg.Offset, msg); err != nil {
			return err
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

type MConsumerGroupConfig struct {
	KafkaVersion   sarama.KafkaVersion
	OffsetsInitial int64
	IsReturnErr    bool
}

func NewMConsumerGroup(consumerConfig *MConsumerGroupConfig, topics, addr []string, groupID string) *MConsumerGroup {
	config := sarama.NewConfig()
	config.Version = consumerConfig.KafkaVersion
	config.Consumer.Offsets.Initial = consumerConfig.OffsetsInitial
	config.Consumer.Return.Errors = consumerConfig.IsReturnErr
	config.Consumer.Offsets.Retry.Max = 99
	config.Consumer.Offsets.AutoCommit.Enable = true
	client, err := sarama.NewClient(addr, config)
	if err != nil {
		panic(err.Error())
	}
	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		panic(err.Error())
	}
	return &MConsumerGroup{
		consumerGroup,
		groupID,
		topics,
	}
}
func (mc *MConsumerGroup) RegisterHandleAndConsumer(handler msgConsumer) {
	for {
		err := mc.Consume(context.Background(),
			mc.topics, consumer{handler})
		if err != nil {
			panic(err.Error())
		}
	}
}
