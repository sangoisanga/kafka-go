package franz

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"kafka/pkg/interfaces"
	"kafka/pkg/services/logger"
)

func NewClient(groups []*Group) interfaces.Client {
	return &client{groups: groups}
}

type client struct {
	groups []*Group
}

func (c *client) StartConsumers() {
	for _, group := range c.groups {
		fmt.Println(group)

		readerConfig := kafka.ReaderConfig{
			Brokers: group.Brokers,
			GroupID: group.GroupID,
			Topic:   group.Topic,
		}

		logger.I.Info(fmt.Sprintf("start %d consumer for topic %s", group.NumOfConsumer, group.Topic))
		for i := 0; i < group.NumOfConsumer; i++ {
			reader := kafka.NewReader(readerConfig)
			go consumeMessage(group, reader)
		}
	}
}

func consumeMessage(group *Group, reader *kafka.Reader) {
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			logger.I.Error("fail when get message", zap.Error(err))
		}
	}(reader)

	ctx := context.Background()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.I.Error("fail when get message", zap.Error(err))
		}

		retry, err := group.ConsumeHandler.Consume(ctx, message)

		if err != nil {
			logger.I.Error("fail when consume message", zap.Error(err))
		}

		if retry {
			err = group.RetryHandler.Retry(ctx, message)
			if err != nil {
				logger.I.Error("retry fail", zap.Error(err))
			}
		}
	}
}
