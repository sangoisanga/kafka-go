package repositories

import (
	"context"
	"github.com/sangoisanga/kafka-go/pkg/interfaces"
	"github.com/segmentio/kafka-go"
)

func BackLogRepository(dbInstance Database, collectionName string) interfaces.BackLogRepository {
	return &backLogRepository{dbInstance: dbInstance, collectionName: collectionName}
}

type backLogRepository struct {
	dbInstance     Database
	collectionName string
}

func (b backLogRepository) Insert(ctx context.Context, message kafka.Message) error {
	coll := b.dbInstance.Instance().Collection(b.collectionName)
	_, err := coll.InsertOne(ctx, message)
	return err
}
