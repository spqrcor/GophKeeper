package application

import (
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/models"
	"GophKeeper/internal/client/transport"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"time"
)

// formatTitle краткое отображение записи
func (app *Application) formatTitle(data models.ItemData) string {
	itemType := app.getTypeDescription(data.Type)
	switch data.Type {
	case "TEXT":
		return fmt.Sprintf("%s: *****", itemType)
	case "CARD":
		return fmt.Sprintf("%s: %s", itemType, data.CardNum)
	case "FILE":
		return fmt.Sprintf("%s: %s", itemType, data.FileName)
	case "AUTH":
		return fmt.Sprintf("%s: %s/****", itemType, data.Login)
	}
	return app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "unknownType",
		},
	})
}

// formatFullText описание текста
func (app *Application) formatFullText(data models.ItemData) string {
	text := app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "type",
		},
	}) + ": " + app.getTypeDescription(data.Type) + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "added",
		},
	}) + ": " + formatDescription(data) + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "text",
		},
	}) + ": " + data.Text
	return text
}

// formatFullCard описание карты
func (app *Application) formatFullCard(data models.ItemData) string {
	text := app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "type",
		},
	}) + ": " + app.getTypeDescription(data.Type) + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "added",
		},
	}) + ": " + formatDescription(data) + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "cardNumber",
		},
	}) + ": " + data.CardNum + "\n"
	text += "CVV: " + data.CardPin + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "validDate",
		},
	}) + ": " + data.CardValid + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "payer",
		},
	}) + ": " + data.CardPayer
	return text
}

// formatFullFile описание файла
func (app *Application) formatFullFile(data models.ItemData, config config.Config, transportData *transport.Data) string {
	text := app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "type",
		},
	}) + ": " + app.getTypeDescription(data.Type) + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "added",
		},
	}) + ": " + formatDescription(data) + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "name",
		},
	}) + ": " + data.FileName + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "link",
		},
	}) + ": " + config.Api + "/api/items/file/" + data.Id + "/token/" + transportData.Token
	return text
}

// formatFullTextAuth описание авторизации
func (app *Application) formatFullTextAuth(data models.ItemData) string {
	text := app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "type",
		},
	}) + ": " + app.getTypeDescription(data.Type) + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "added",
		},
	}) + ": " + formatDescription(data) + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "login",
		},
	}) + ": " + data.Login + "\n"
	text += app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "password",
		},
	}) + ": " + data.Password
	return text
}

// formatFull полное отображение записи
func (app *Application) formatFull(data models.ItemData, config config.Config, transportData *transport.Data) string {
	switch data.Type {
	case "TEXT":
		return app.formatFullText(data)
	case "CARD":
		return app.formatFullCard(data)
	case "FILE":
		return app.formatFullFile(data, config, transportData)
	case "AUTH":
		return app.formatFullTextAuth(data)
	}
	return app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "unknownType",
		},
	})
}

// getTypeDescription получение описания типа
func (app *Application) getTypeDescription(itemType string) string {
	switch itemType {
	case "TEXT":
		return app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "text",
			},
		})
	case "CARD":
		return app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "card",
			},
		})
	case "FILE":
		return app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "file",
			},
		})
	case "AUTH":
		return app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "authorization",
			},
		})
	}
	return app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "unknownType",
		},
	})
}

// formatDescription форматирование даты
func formatDescription(data models.ItemData) string {
	date, _ := time.Parse(time.RFC3339, data.CreatedAt)
	return date.Format("15:04 02.01.2006")
}
