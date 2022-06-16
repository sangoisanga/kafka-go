package franz

import (
	"fmt"
	"kafka/pkg/franz/handler"
	"kafka/pkg/franz/retry"
	"kafka/pkg/interfaces"
	"kafka/pkg/models"
	"kafka/pkg/repositories"
)

func Factory(db repositories.Database) interfaces.Factory {
	return &factory{dbInstance: db}
}

type factory struct {
	dbInstance repositories.Database
}

const DeadLetterPostfix = "dead-letter"

func (f *factory) GetConsumerHandler(topic string) interfaces.ConsumerHandler {
	switch topic {
	case "topic1":
		return handler.PrintLogHandler()
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
		return retry.DeadLetterStrategy(producer)

	case "back-log":
		repo := repositories.BackLogRepository(f.dbInstance, config.BackLogCollection)
		return retry.BackLogStrategy(repo)

	default:
		panic("undefined strategy")
	}
}
