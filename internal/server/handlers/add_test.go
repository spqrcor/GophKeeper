package handlers

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/mocks"
	"GophKeeper/internal/server/storage"
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddItemHandler(t *testing.T) {
	conf := config.NewConfig()
	tokenAuth := jwtauth.New("HS256", []byte(conf.SecretKey), nil, jwt.WithAcceptableSkew(conf.TokenExp))
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	store := mocks.NewMockStorage(mockCtrl)

	store.EXPECT().AddItem(context.Background(), storage.CommonData{Type: "TEXT", Text: "1"}, "672124b6-9894-11e5-be38-001d42e813fe", "1234", nil).Return("", pgx.ErrNoRows).AnyTimes()
	store.EXPECT().AddItem(context.Background(), storage.CommonData{Type: "TEXT", Text: "1"}, "772124b6-9894-11e5-be38-001d42e813fe", "1234", nil).Return("772124b6-9894-11e5-be38-001d42e813fe", nil).AnyTimes()

	tests := []struct {
		name        string
		userID      string
		contentType string
		body        []byte
		statusCode  int
	}{
		{
			name:        "error json",
			userID:      "172124b6-9894-11e5-be38-001d42e813fe",
			contentType: "application/json",
			body:        []byte(`333`),
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "error valid",
			userID:      "172124b6-9894-11e5-be38-001d42e813fe",
			contentType: "application/json",
			body:        []byte(`{"type":"TEXT"}`),
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "http 500",
			userID:      "672124b6-9894-11e5-be38-001d42e813fe",
			contentType: "application/json",
			body:        []byte(`{"type":"TEXT","TEXT":"1"}`),
			statusCode:  http.StatusInternalServerError,
		},
		{
			name:        "success simple",
			userID:      "772124b6-9894-11e5-be38-001d42e813fe",
			contentType: "application/json",
			body:        []byte(`{"type":"TEXT","TEXT":"1"}`),
			statusCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"UserID": tt.userID, "PIN": "1234"})
			r := chi.NewRouter()
			r.Use(
				jwtauth.Verifier(tokenAuth),
				jwtauth.Authenticator(tokenAuth),
			)

			r.Post("/api/items", AddItemHandler(store))
			ts := httptest.NewServer(r)
			defer ts.Close()
			req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/items", bytes.NewReader(tt.body))

			req.Header.Add("Content-Type", tt.contentType)
			req.Header.Set("Authorization", "Bearer "+tokenString)
			resp, _ := http.DefaultClient.Do(req)
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
		})
	}

}
