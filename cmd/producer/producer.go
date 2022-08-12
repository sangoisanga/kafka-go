package main

import (
	"context"
	"fmt"
	"github.com/sangoisanga/core-go/pkg/log"
	"github.com/sangoisanga/kafka-go/pkg/franz"
	"github.com/sangoisanga/kafka-go/pkg/interfaces"
	"github.com/sangoisanga/kafka-go/pkg/services/config"
	"go.uber.org/zap"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func main() {
	log.InitLogger()
	logger := log.Logger()
	configs := config.Load()
	producer := franz.NewProducer(franz.ProducerConfig{
		Address: configs.KafkaURL,
		Topic:   configs.Topic,
	})

	defer func(producer interfaces.Producer) {
		err := producer.Close()
		if err != nil {
			logger.Error("closing writer error", zap.Error(err))
		}
	}(producer)

	logger.Info("start producing ... !!")
	for i := 0; ; i++ {
		key := fmt.Sprintf("Key-%d", i)
		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte(fmt.Sprint(uuid.New())),
		}
		err := producer.Write(context.Background(), msg)
		if err != nil {
			logger.Error("write error", zap.Error(err))
		} else {
			logger.Info(fmt.Sprintf("produced %s", key))
		}
		time.Sleep(1 * time.Second)
	}
}
