package franz

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"kafka/pkg/interfaces"
	"kafka/pkg/services/logger"
	"time"
)

func DeadLetterStrategy(producer Producer) interfaces.RetryStrategy {
	return &deadLetter{producer: producer}
}

type deadLetter struct {
	producer Producer
}

func (d deadLetter) Retry(ctx context.Context, message kafka.Message) error {
	var err error
	retries := 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		message.Topic = "" // override old topic to dead letter topic
		err = d.producer.Write(ctx, message)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}
		if err != nil {
			logger.I.Error("write fail", zap.Error(err))
		}
	}
	return err
}
