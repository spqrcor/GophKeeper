package logger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger(zap.InfoLevel)
	assert.Nil(t, err)
	assert.NotNil(t, logger)
}
