package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"kafka/pkg/services/config"
	"kafka/pkg/services/logger"
	"time"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func main() {
	logger.Init()
	configs := config.Load()
	writer := newKafkaWriter(configs.KafkaURL, configs.Topic)

	defer func(writer *kafka.Writer) {
		err := writer.Close()
		if err != nil {
			logger.I.Error("closing writer error", zap.Error(err))
		}
	}(writer)

	logger.I.Info("start producing ... !!")
	for i := 0; ; i++ {
		key := fmt.Sprintf("Key-%d", i)
		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte(fmt.Sprint(uuid.New())),
		}
		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			logger.I.Error("write error", zap.Error(err))
		} else {
			logger.I.Info(fmt.Sprintf("produced %s", key))
		}
		time.Sleep(1 * time.Second)
	}
}
