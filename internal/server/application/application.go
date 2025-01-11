// Package application https server
package application

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/handlers"
	"GophKeeper/internal/server/storage/db"
	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Application тип Application
type Application struct {
	config    config.Config
	logger    *zap.Logger
	tokenAuth *jwtauth.JWTAuth
	dbRes     *sql.DB
}

// NewApplication создание Application, opts набор параметров
func NewApplication(opts ...func(*Application)) *Application {
	app := &Application{}
	for _, opt := range opts {
		opt(app)
	}
	return app
}

// WithLogger добавление logger
func WithLogger(logger *zap.Logger) func(*Application) {
	return func(a *Application) {
		a.logger = logger
	}
}

// WithConfig добавление config
func WithConfig(config config.Config) func(*Application) {
	return func(a *Application) {
		a.config = config
	}
}

// WithTokenAuth добавление TokenAuth
func WithTokenAuth(tokenAuth *jwtauth.JWTAuth) func(*Application) {
	return func(a *Application) {
		a.tokenAuth = tokenAuth
	}
}

// WithDB добавление storage
func WithDB(dbRes *sql.DB) func(*Application) {
	return func(a *Application) {
		a.dbRes = dbRes
	}
}

// NewHTTPServer создание http сервера
func (a *Application) NewHTTPServer() *http.Server {
	r := chi.NewRouter()
	r.Use(LoggerMiddleware(a.logger))
	r.Use(middleware.Compress(5, "application/json", "text/plain"))

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(a.tokenAuth))
		r.Use(jwtauth.Authenticator(a.tokenAuth))
		r.Get("/api/items", handlers.GetItemsHandler(db.CreateListItemDB(a.config, a.logger, a.dbRes)))
		r.Delete("/api/items/{id}", handlers.RemoveItemHandler(db.CreateRemoveItemDB(a.config, a.dbRes)))
		r.Post("/api/items", handlers.AddItemHandler(db.CreateAddItemDB(a.config, a.dbRes)))
		r.Post("/api/items/file", handlers.AddItemFileHandler(db.CreateAddItemDB(a.config, a.dbRes), a.config.MaxUploadFileSize))
	})

	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", handlers.RegisterHandler(db.CreateRegisterUserDB(a.config, a.dbRes)))
		r.Post("/api/user/login", handlers.LoginHandler(db.CreateLoginUserDB(a.config, a.dbRes), a.tokenAuth))
		r.Get("/api/items/file/{id}/token/{token}", handlers.GetItemFileHandler(db.CreateItemInfoDB(a.config, a.logger, a.dbRes), a.tokenAuth))
	})

	r.HandleFunc(`/*`, func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
	})

	return &http.Server{
		Handler: r,
		Addr:    a.config.Addr,
	}
}

// Start запуск приложения
func (a *Application) Start() {
	httpServer := a.NewHTTPServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-stop
		if err := httpServer.Shutdown(context.Background()); err != nil {
			a.logger.Error(err.Error())
		}
		if err := a.dbRes.Close(); err != nil {
			a.logger.Error(err.Error())
		}
	}()

	if err := initCertificate(); err != nil {
		a.logger.Error(err.Error())
	}
	if err := httpServer.ListenAndServeTLS(certCfg.certPath, certCfg.keyPath); err != nil {
		a.logger.Error(err.Error())
	}
	a.logger.Info("graceful shutdown")
}
