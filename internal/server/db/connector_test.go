package db

import (
	"GophKeeper/internal/server/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnect(t *testing.T) {
	conf := config.NewConfig()
	if conf.DatabaseDSN == "" {
		t.Skip("Skipping testing...")
	}

	tests := []struct {
		name string
		dsn  string
		want bool
	}{
		{
			"Error",
			"",
			false,
		},
		{
			"Success",
			conf.DatabaseDSN,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Connect(tt.dsn)
			assert.Equal(t, tt.want, err == nil)
		})
	}
}
