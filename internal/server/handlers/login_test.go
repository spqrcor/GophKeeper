package handlers

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/mocks"
	"GophKeeper/internal/server/storage"
	"bytes"
	"context"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	conf := config.NewConfig()
	authenticate := jwtauth.New("HS256", []byte(conf.SecretKey), nil, jwt.WithAcceptableSkew(conf.TokenExp))
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	store := mocks.NewMockStorage(mockCtrl)

	store.EXPECT().Login(context.Background(), storage.InputDataUser{
		Login: "spqr", Password: "1",
	}).Return("", storage.ErrLogin).AnyTimes()
	store.EXPECT().Login(context.Background(), storage.InputDataUser{
		Login: "spqr", Password: "123456", Pin: "123456",
	}).Return("672124b6-9894-11e5-be38-001d42e813fe", nil).AnyTimes()
	store.EXPECT().Login(context.Background(), storage.InputDataUser{
		Login: "xxx", Password: "xxx", Pin: "123456",
	}).Return("", pgx.ErrNoRows).AnyTimes()
	store.EXPECT().Login(context.Background(), storage.InputDataUser{
		Login: "xxx2", Password: "xxx2", Pin: "",
	}).Return("672124b6-9894-11e5-be38-001d42e813fe", nil).AnyTimes()

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
			name:        "login error",
			contentType: "application/json",
			body:        []byte(`{"login":"spqr","password":"1"}`),
			statusCode:  http.StatusUnauthorized,
		},
		{
			name:        "success",
			contentType: "application/json",
			body:        []byte(`{"login":"spqr","password":"123456","pin":"123456"}`),
			statusCode:  http.StatusOK,
		},
		{
			name:        "http 500",
			contentType: "application/json",
			body:        []byte(`{"login":"xxx","password":"xxx","pin":"123456"}`),
			statusCode:  http.StatusInternalServerError,
		},
		{
			name:        "no pin error",
			contentType: "application/json",
			body:        []byte(`{"login":"xxx2","password":"xxx2","pin":""}`),
			statusCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(tt.body))
			req.Header.Add("Content-Type", tt.contentType)
			LoginHandler(store, authenticate)(rw, req)

			resp := rw.Result()
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
			_ = resp.Body.Close()
		})
	}
}
