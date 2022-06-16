package main

import (
	"kafka/pkg/franz"
	"kafka/pkg/repositories"
	"kafka/pkg/services/config"
	"kafka/pkg/services/logger"
	"os"
	"os/signal"
)

func main() {
	logger.Init()
	configs := config.Load()

	mongoClient, err := repositories.CreateNewClient(configs.MongodbURI)
	if err != nil {
		panic(err)
	}

	backLogDB := repositories.Database{Client: mongoClient, DBName: configs.BackLogDB.DBName}
	franzFactory := franz.Factory(backLogDB)

	logger.I.Info("init consumer")
	groups := make([]*franz.Group, 0, len(configs.KafkaConfig))
	for _, groupConfig := range configs.KafkaConfig {
		group := &franz.Group{
			Brokers:        groupConfig.Brokers,
			Topic:          groupConfig.Topic,
			GroupID:        groupConfig.GroupID,
			NumOfConsumer:  groupConfig.NumOfConsumer,
			ConsumeHandler: franzFactory.GetConsumerHandler(groupConfig.Topic),
			RetryHandler:   franzFactory.GetRetryStrategy(groupConfig),
		}
		groups = append(groups, group)
	}

	client := franz.NewClient(groups)
	client.StartConsumers()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.I.Info("exiting ...")
}
