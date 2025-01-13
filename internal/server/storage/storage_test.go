package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserValidator(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		password string
		result   bool
	}{
		{
			name:     "error password",
			login:    "spqr",
			password: "1",
			result:   false,
		},
		{
			name:     "error login",
			login:    "1",
			password: "123456",
			result:   false,
		},
		{
			name:     "success",
			login:    "spqr",
			password: "123456",
			result:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := UserValidator(InputDataUser{Login: tt.login, Password: tt.password})
			assert.Equal(t, tt.result, res == nil)
		})
	}
}

func TestItemValidator(t *testing.T) {
	tests := []struct {
		name   string
		item   CommonData
		result bool
	}{
		{
			name:   "error1",
			item:   CommonData{Type: "---"},
			result: false,
		},
		{
			name:   "error2",
			item:   CommonData{Type: "TEXT"},
			result: false,
		},
		{
			name:   "error3",
			item:   CommonData{Type: "CARD"},
			result: false,
		},
		{
			name:   "error4",
			item:   CommonData{Type: "AUTH"},
			result: false,
		},
		{
			name:   "error5",
			item:   CommonData{Type: "FILE"},
			result: false,
		},
		{
			name:   "success",
			item:   CommonData{Type: "FILE", FileName: "test.txt"},
			result: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ItemValidator(tt.item)
			assert.Equal(t, tt.result, res == nil)
		})
	}
}
