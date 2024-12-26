package main

import (
	"GophKeeper/internal/client/application"
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/logger"
	"GophKeeper/internal/client/transport"
	"log"
)

func main() {
	conf := config.NewConfig()
	loggerRes, err := logger.NewLogger(conf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	trans := transport.NewTransport(conf, loggerRes)

	app := application.NewApplication(
		application.WithConfig(conf),
		application.WithLogger(loggerRes),
		application.WithTransport(trans),
	)
	app.Start()

}
