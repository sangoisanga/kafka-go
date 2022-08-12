package retry

import (
	"context"
	"errors"
	"github.com/sangoisanga/core-go/pkg/log"
	"github.com/sangoisanga/kafka-go/pkg/interfaces"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

func DeadLetterStrategy(producer interfaces.Producer) interfaces.RetryStrategy {
	return &deadLetter{producer: producer}
}

type deadLetter struct {
	producer interfaces.Producer
}

func (d deadLetter) Retry(ctx context.Context, message kafka.Message) error {
	logger := log.Logger()
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
			logger.Error("write fail", zap.Error(err))
		}
	}
	return err
}
