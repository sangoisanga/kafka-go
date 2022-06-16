package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"kafka/pkg/franz"
	"kafka/pkg/interfaces"
	"kafka/pkg/services/config"
	"kafka/pkg/services/logger"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func main() {
	logger.Init()
	configs := config.Load()
	producer := franz.NewProducer(franz.ProducerConfig{
		Address: configs.KafkaURL,
		Topic:   configs.Topic,
	})

	defer func(producer interfaces.Producer) {
		err := producer.Close()
		if err != nil {
			logger.I.Error("closing writer error", zap.Error(err))
		}
	}(producer)

	logger.I.Info("start producing ... !!")
	for i := 0; ; i++ {
		key := fmt.Sprintf("Key-%d", i)
		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte(fmt.Sprint(uuid.New())),
		}
		err := producer.Write(context.Background(), msg)
		if err != nil {
			logger.I.Error("write error", zap.Error(err))
		} else {
			logger.I.Info(fmt.Sprintf("produced %s", key))
		}
		time.Sleep(1 * time.Second)
	}
}
