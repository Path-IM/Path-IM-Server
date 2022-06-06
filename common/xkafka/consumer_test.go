package xkafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"math"
	"testing"
)

type consumer struct {
}

func (c consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("%s", string(msg.Value))
		//sess.MarkMessage(msg, "")
	}
	return nil
}

func TestConsume(t *testing.T) {
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_2_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Retry.Max = math.MaxInt
	config.Consumer.Offsets.AutoCommit.Enable = true
	client, err := sarama.NewClient([]string{
		"172.27.10.3:9092",
	}, config)
	if err != nil {
		panic(err.Error())
	}
	consumerGroup, err := sarama.NewConsumerGroupFromClient("test_msg", client)
	if err != nil {
		panic(err.Error())
	}
	for {
		err := consumerGroup.Consume(context.Background(), []string{
			"im_msg",
		}, consumer{})
		if err != nil {
			panic(err.Error())
		}
	}
}
