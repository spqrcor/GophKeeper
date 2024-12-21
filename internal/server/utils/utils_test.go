package utils

import (
	"GophKeeper/internal/server/config"
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestCreateKeyFromPin(t *testing.T) {
	conf := config.NewConfig()
	tests := []struct {
		name   string
		pin    string
		hash   string
		result bool
	}{
		{
			name:   "success",
			pin:    "1234",
			hash:   "f0a35cde7acd30194ab8417b067097fe79c290226a9f7e2ae358512e908cd057",
			result: true,
		},
		{
			name:   "error",
			pin:    "12345",
			hash:   "10a35cde7acd30194ab8417b067097fe79c290226a9f7e2ae358512e908cd057",
			result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, CreateKeyFromPin(tt.pin, conf.Salt) == tt.hash, tt.result)
		})
	}
}

func TestFromPostJSON(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        []byte
		result      bool
	}{
		{
			name:        "success",
			contentType: "application/json",
			body:        []byte(`{"login":"xxx2","password":"xxx2","pin":""}`),
			result:      true,
		},
		{
			name:        "error",
			contentType: "plain/text",
			body:        []byte(`55`),
			result:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(tt.body))
			req.Header.Add("Content-Type", tt.contentType)
			input := ""
			res := FromPostJSON(req, input)
			assert.Equal(t, tt.result, res == nil)
		})
	}
}
