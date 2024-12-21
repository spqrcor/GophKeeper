package application

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/handlers"
	"GophKeeper/internal/server/logger"
	"GophKeeper/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"

	"testing"
)

func Test_loggerMiddleware(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	store := storage.NewStorage(conf, loggerRes)

	r := chi.NewRouter()
	r.Use(LoggerMiddleware(loggerRes))

	r.Post(`/`, handlers.RegisterHandler(store))
	srv := httptest.NewServer(r)
	defer srv.Close()

	t.Run("work_logger", func(t *testing.T) {
		r := httptest.NewRequest("POST", srv.URL+"/", nil)
		r.RequestURI = ""
		r.Header.Set("Content-Type", "text/html")
		resp, _ := http.DefaultClient.Do(r)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		defer func() {
			err := resp.Body.Close()
			require.NoError(t, err)
		}()
	})
}
