package xkafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/showurl/Zero-IM-Server/common/xtrace"
	"go.opentelemetry.io/otel/attribute"
)

type Producer struct {
	topic    string
	addr     []string
	config   *sarama.Config
	producer sarama.SyncProducer
}

func MustNewProducer(config ProducerConfig) *Producer {
	p := Producer{}
	p.config = sarama.NewConfig()                             //Instantiate a sarama Config
	p.config.Producer.Return.Successes = true                 //Whether to enable the successes channel to be notified after the message is sent successfully
	p.config.Producer.RequiredAcks = sarama.WaitForAll        //Set producer Message Reply level 0 1 all
	p.config.Producer.Partitioner = sarama.NewHashPartitioner //Set the hash-key automatic hash partition. When sending a message, you must specify the key value of the message. If there is no key, the partition will be selected randomly

	p.addr = config.Brokers
	p.topic = config.Topic

	producer, err := sarama.NewSyncProducer(p.addr, p.config) //Initialize the client
	if err != nil {
		panic(err.Error())
	}
	p.producer = producer
	return &p
}

func (p *Producer) SendMessage(ctx context.Context, m proto.Message, key ...string) (partition int32, offset int64, err error) {
	kMsg := &sarama.ProducerMessage{}
	kMsg.Topic = p.topic
	if len(key) == 1 {
		kMsg.Key = sarama.StringEncoder(key[0])
	}
	bMsg, err := proto.Marshal(m)
	if err != nil {
		return -1, -1, err
	}
	kMsg.Value = sarama.ByteEncoder(bMsg)
	xtrace.StartFuncSpan(ctx, "SendMessageToKafka", func(ctx context.Context) {
		partition, offset, err = p.producer.SendMessage(kMsg)
	},
		attribute.StringSlice("keys", key),
		attribute.String("topic", p.topic),
	)
	return
}
