package handlers

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/mocks"
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

func TestRemoveItemHandler(t *testing.T) {
	conf := config.NewConfig()
	tokenAuth := jwtauth.New("HS256", []byte(conf.SecretKey), nil, jwt.WithAcceptableSkew(conf.TokenExp))
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	store := mocks.NewMockStorage(mockCtrl)

	store.EXPECT().RemoveItem(context.Background(), "672124b6-9894-11e5-be38-001d42e813fe", "672124b6-9894-11e5-be38-001d42e813fe").Return(pgx.ErrNoRows).AnyTimes()
	store.EXPECT().RemoveItem(context.Background(), "872124b6-9894-11e5-be38-001d42e813fe", "872124b6-9894-11e5-be38-001d42e813fe").Return(nil).AnyTimes()

	tests := []struct {
		name       string
		userID     string
		itemID     string
		statusCode int
	}{
		{
			name:       "http 500",
			userID:     "672124b6-9894-11e5-be38-001d42e813fe",
			itemID:     "672124b6-9894-11e5-be38-001d42e813fe",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "http 200",
			userID:     "872124b6-9894-11e5-be38-001d42e813fe",
			itemID:     "872124b6-9894-11e5-be38-001d42e813fe",
			statusCode: http.StatusOK,
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

			r.Delete("/api/items/{id}", RemoveItemHandler(store))
			ts := httptest.NewServer(r)
			defer ts.Close()
			req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/items/"+tt.itemID, nil)

			req.Header.Set("Authorization", "Bearer "+tokenString)
			resp, _ := http.DefaultClient.Do(req)
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
		})
	}

}
