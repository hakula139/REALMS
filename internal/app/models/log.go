package models

import (
	"go.uber.org/zap"
)

// InitLogger initializes a logger to record all events and write to log files
func InitLogger(cfg zap.Config) (*zap.Logger, error) {
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	// logger.Info("Logger construction succeeded")
	return logger, nil
}
