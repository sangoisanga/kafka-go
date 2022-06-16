package franz

import (
	"fmt"
	"kafka/pkg/interfaces"
	"kafka/pkg/models"
	"kafka/pkg/repositories"
)

func Factory(db repositories.Database) *factory {
	return &factory{dbInstance: db}
}

type factory struct {
	dbInstance repositories.Database
}

const DeadLetterPostfix = "dead-letter"

func (f *factory) GetConsumerHandler(topic string) interfaces.ConsumerHandler {
	switch topic {
	case "topic1":
		return PrintLogHandler()
	default:
		panic("undefined handler")
	}
}

func (f *factory) GetRetryStrategy(config *models.KafkaGroupConfig) interfaces.RetryStrategy {
	switch config.RetryStrategy {
	case "dead-letter":
		topic := fmt.Sprintf("%s.%s", config.Topic, DeadLetterPostfix)
		producerConfig := ProducerConfig{
			Address: config.DeadLetterAddress,
			Topic:   topic,
		}
		producer := NewProducer(producerConfig)
		return DeadLetterStrategy(producer)

	case "back-log":
		repo := repositories.BackLogRepository(f.dbInstance, config.BackLogCollection)
		return BackLogStrategy(repo)

	default:
		panic("undefined strategy")
	}
}
