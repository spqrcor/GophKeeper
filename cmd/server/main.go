package main

import (
	"GophKeeper/internal/server/application"
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/logger"
	"GophKeeper/internal/server/storage"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"log"
)

func main() {
	conf := config.NewConfig()
	loggerRes, err := logger.NewLogger(conf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	storageService := storage.NewStorage(conf, loggerRes)
	app := application.NewApplication(
		application.WithConfig(conf),
		application.WithLogger(loggerRes),
		application.WithStorage(storageService),
		application.WithTokenAuth(jwtauth.New("HS256", []byte(conf.SecretKey), nil, jwt.WithAcceptableSkew(conf.TokenExp))),
	)
	app.Start()
}
