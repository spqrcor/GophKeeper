package main

import (
	"GophKeeper/internal/server/application"
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/db"
	"GophKeeper/internal/server/logger"
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

	dbres, err := db.Connect(conf.DatabaseDSN)
	if err != nil {
		loggerRes.Fatal(err.Error())
	}
	if err := db.Migrate(dbres); err != nil {
		loggerRes.Fatal(err.Error())
	}

	app := application.NewApplication(
		application.WithConfig(conf),
		application.WithLogger(loggerRes),
		application.WithTokenAuth(jwtauth.New("HS256", []byte(conf.SecretKey), nil, jwt.WithAcceptableSkew(conf.TokenExp))),
		application.WithDB(dbres),
	)
	app.Start()
}
