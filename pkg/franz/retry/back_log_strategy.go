package retry

import (
	"context"
	"github.com/sangoisanga/kafka-go/interfaces"
	"github.com/segmentio/kafka-go"
)

func BackLogStrategy(repo interfaces.BackLogRepository) interfaces.RetryStrategy {
	return &backLogStrategy{repo: repo}
}

type backLogStrategy struct {
	repo interfaces.BackLogRepository
}

func (b backLogStrategy) Retry(ctx context.Context, message kafka.Message) error {
	return b.repo.Insert(ctx, message)
}
