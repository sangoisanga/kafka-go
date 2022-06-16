package franz

import (
	"context"
	"github.com/segmentio/kafka-go"
	"kafka/pkg/interfaces"
)

type ProducerConfig struct {
	Address string
	Topic   string
}

func NewProducer(config ProducerConfig) interfaces.Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(config.Address),
		Topic:                  config.Topic,
		AllowAutoTopicCreation: true,
	}
	return &producer{writer: w}

}

type producer struct {
	writer *kafka.Writer
}

func (p *producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

func (p *producer) Write(ctx context.Context, message kafka.Message) error {
	return p.writer.WriteMessages(ctx, message)
}
