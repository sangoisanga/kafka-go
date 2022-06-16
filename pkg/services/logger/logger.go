package logger

import (
	"go.uber.org/zap"
	"sync"
)

var I *zap.Logger
var once = sync.Once{}

func Init() {
	once.Do(func() {
		I, _ = zap.NewProduction()
	})
}
