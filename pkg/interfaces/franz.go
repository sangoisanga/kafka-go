package interfaces

import (
	"context"
	"github.com/sangoisanga/kafka-go/pkg/models"
	"github.com/segmentio/kafka-go"
)

type ConsumerHandler interface {
	Consume(ctx context.Context, message kafka.Message) (retry bool, err error)
}

type RetryStrategy interface {
	Retry(ctx context.Context, message kafka.Message) error
}

type BackLogRepository interface {
	Insert(ctx context.Context, message kafka.Message) error
}

type Producer interface {
	Write(ctx context.Context, message kafka.Message) error
	Close() error
}

type Factory interface {
	GetConsumerHandler(topic string) ConsumerHandler
	GetRetryStrategy(config *models.KafkaGroupConfig) RetryStrategy
}

type Client interface {
	StartConsumers()
}
