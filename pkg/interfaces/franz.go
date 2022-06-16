package interfaces

import (
	"context"
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
