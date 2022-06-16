package franz

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
	Address string
	Topic   string
}

type Producer interface {
	Write(ctx context.Context, message kafka.Message) error
	Close()
}

func NewProducer(config ProducerConfig) Producer {
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

func (p *producer) Close() {
	if p.writer != nil {
		p.writer.Close()
	}
}

func (p *producer) Write(ctx context.Context, message kafka.Message) error {
	return p.writer.WriteMessages(ctx, message)
}
