package main

import (
	"context"
	kafka "github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"kafka/pkg/franz"
	"kafka/pkg/repositories"
	"kafka/pkg/services/config"
	"kafka/pkg/services/logger"
	"os"
	"os/signal"
)

func getMongoCollection(mongoURL, dbName, collectionName string) *mongo.Collection {
	option := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context.Background(), option)
	if err != nil {
		logger.I.Error("connect db error", zap.Error(err))
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		logger.I.Error("ping error", zap.Error(err))
	}

	logger.I.Info("Connected to MongoDB ... !!")

	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	return collection
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

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
