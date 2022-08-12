package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sangoisanga/core-go/pkg/log"
	"github.com/sangoisanga/kafka-go/interfaces"
	"github.com/segmentio/kafka-go"
	"time"
)

func PrintLogHandler() interfaces.ConsumerHandler {
	return &printLogHandler{}
}

type printLogHandler struct {
}

func (p printLogHandler) Consume(_ context.Context, message kafka.Message) (retry bool, err error) {
	logger := log.Logger()
	bytes, err := json.Marshal(message)
	if err != nil {
		return true, err
	}
	logger.Info(string(bytes))

	if time.Now().Second()%2 == 0 {
		return true, errors.New("error is real")
	}

	return false, nil
}
