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
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestAddItemFileHandler(t *testing.T) {
	conf := config.NewConfig()
	tokenAuth := jwtauth.New("HS256", []byte(conf.SecretKey), nil, jwt.WithAcceptableSkew(conf.TokenExp))
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	store := mocks.NewMockStorage(mockCtrl)

	store.EXPECT().AddItem(context.Background(), storage.CommonData{Type: "FILE", FileName: "test1.txt"}, "672124b6-9894-11e5-be38-001d42e813fe", "1234", []byte("1")).Return("672124b6-9894-11e5-be38-001d42e813fe", nil).AnyTimes()
	store.EXPECT().AddItem(context.Background(), storage.CommonData{Type: "FILE", FileName: "test2.txt"}, "872124b6-9894-11e5-be38-001d42e813fe", "1234", []byte("2")).Return("", pgx.ErrNoRows).AnyTimes()

	tests := []struct {
		name       string
		userID     string
		fileName   string
		fileKey    string
		statusCode int
	}{
		{
			name:       "http 200",
			userID:     "672124b6-9894-11e5-be38-001d42e813fe",
			fileName:   "test1.txt",
			fileKey:    "file",
			statusCode: http.StatusOK,
		},
		{
			name:       "http 400",
			userID:     "672124b6-9894-11e5-be38-001d42e813fe",
			fileName:   "test2.txt",
			fileKey:    "file2",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "http 500",
			userID:     "872124b6-9894-11e5-be38-001d42e813fe",
			fileName:   "test2.txt",
			fileKey:    "file",
			statusCode: http.StatusInternalServerError,
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
			r.Post("/api/items/file", AddItemFileHandler(store))
			ts := httptest.NewServer(r)
			defer ts.Close()

			file, _ := os.Open("../../../data/test/" + tt.fileName)
			defer file.Close()
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, _ := writer.CreateFormFile(tt.fileKey, filepath.Base(tt.fileName))
			_, _ = io.Copy(part, file)
			_ = writer.Close()

			req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/items/file", body)
			req.Header.Add("Content-Type", writer.FormDataContentType())
			req.Header.Set("Authorization", "Bearer "+tokenString)
			resp, _ := http.DefaultClient.Do(req)
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
		})
	}

}
