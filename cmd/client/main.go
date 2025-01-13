package main

import (
	"GophKeeper/internal/client/application"
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/logger"
	"GophKeeper/internal/client/transport/rest"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
)

func main() {
	conf := config.NewConfig()
	loggerRes, err := logger.NewLogger(conf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	trans := rest.CreateRestTransport(conf, loggerRes)

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile(conf.LanguagePath + "/" + conf.Language + ".toml")
	localizer := i18n.NewLocalizer(bundle, conf.Language)

	xxx := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "setPIN",
		},
	})
	fmt.Printf(xxx)

	app := application.NewApplication(
		application.WithConfig(conf),
		application.WithLogger(loggerRes),
		application.WithTransport(trans),
		application.WithLocalizer(localizer),
	)
	app.Start()

}
