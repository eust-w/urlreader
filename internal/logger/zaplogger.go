package logger

import (
	"go.uber.org/zap"
	"sync"
)

var (
	logger     *zap.SugaredLogger
	once       sync.Once
)

// InitLogger 初始化全局logger，建议在main中调用
func InitLogger() {
	once.Do(func() {
		l, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		logger = l.Sugar()
	})
}

// GetLogger 获取全局logger
func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		InitLogger()
	}
	return logger
}

// SyncLogger 刷新日志缓冲区，建议在main退出前调用
func SyncLogger() {
	if logger != nil {
		_ = logger.Sync()
	}
}
