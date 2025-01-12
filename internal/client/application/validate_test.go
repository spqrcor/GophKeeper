package application

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_luhnAlgorithm(t *testing.T) {
	tests := []struct {
		name   string
		card   string
		result bool
	}{
		{
			name:   "short card",
			card:   "123456",
			result: false,
		},
		{
			name:   "valid card",
			card:   "5536913839920903",
			result: true,
		},
		{
			name:   "not valid card",
			card:   "5536913839920904",
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, luhnAlgorithm(tt.card), tt.result, "Card validation failed")
		})
	}
}
