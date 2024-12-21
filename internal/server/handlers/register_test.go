package handlers

import (
	"GophKeeper/internal/server/mocks"
	"GophKeeper/internal/server/storage"
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	store := mocks.NewMockStorage(mockCtrl)

	store.EXPECT().Register(context.Background(), storage.InputDataUser{
		Login: "s", Password: "123456",
	}).Return(storage.ErrValidation).AnyTimes()
	store.EXPECT().Register(context.Background(), storage.InputDataUser{
		Login: "spqr", Password: "123456",
	}).Return(nil).AnyTimes().MaxTimes(1)
	store.EXPECT().Register(context.Background(), storage.InputDataUser{
		Login: "spqr", Password: "123456",
	}).Return(storage.ErrLoginExists).AnyTimes().MinTimes(1)
	store.EXPECT().Register(context.Background(), storage.InputDataUser{
		Login: "xxx", Password: "xxx",
	}).Return(pgx.ErrNoRows).AnyTimes()

	tests := []struct {
		name        string
		contentType string
		body        []byte
		statusCode  int
	}{
		{
			name:        "not format error",
			contentType: "application/json",
			body:        []byte(`<num>3333</num>`),
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "validation error",
			contentType: "application/json",
			body:        []byte(`{"login":"s","password":"123456"}`),
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "success",
			contentType: "application/json",
			body:        []byte(`{"login":"spqr","password":"123456"}`),
			statusCode:  http.StatusOK,
		},
		{
			name:        "user exists error",
			contentType: "application/json",
			body:        []byte(`{"login":"spqr","password":"123456"}`),
			statusCode:  http.StatusConflict,
		},
		{
			name:        "http 500",
			contentType: "application/json",
			body:        []byte(`{"login":"xxx","password":"xxx"}`),
			statusCode:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(tt.body))
			req.Header.Add("Content-Type", tt.contentType)
			RegisterHandler(store)(rw, req)

			resp := rw.Result()
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
			_ = resp.Body.Close()
		})
	}
}
