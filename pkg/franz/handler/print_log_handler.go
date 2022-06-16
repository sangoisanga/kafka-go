package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"kafka/pkg/interfaces"
	"kafka/pkg/services/logger"
	"time"
)

func PrintLogHandler() interfaces.ConsumerHandler {
	return &printLogHandler{}
}

type printLogHandler struct {
}

func (p printLogHandler) Consume(_ context.Context, message kafka.Message) (retry bool, err error) {
	bytes, err := json.Marshal(message)
	if err != nil {
		return true, err
	}
	logger.I.Info(string(bytes))

	if time.Now().Second()%2 == 0 {
		return true, errors.New("error is real")
	}

	return false, nil
}
