package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	assert.NotNil(t, config.SecretKey)
	assert.NotNil(t, config.TokenExp)
	assert.NotNil(t, config.LogLevel)
}
