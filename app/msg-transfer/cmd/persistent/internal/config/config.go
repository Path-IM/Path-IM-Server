package config

import (
	"github.com/Path-IM/Path-IM-Server/common/xkafka"
	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	service.ServiceConf
	Kafka KafkaConfig
}
type KafkaConfigOnline struct {
	xkafka.ProducerConfig
	MsgPersistentGroupID string
}
type KafkaConfig struct {
	Online KafkaConfigOnline
}
