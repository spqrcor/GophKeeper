// Package logger логгер
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger инициализация логгера
func NewLogger(logLevel zapcore.Level) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(logLevel)
	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return zl, nil
}
