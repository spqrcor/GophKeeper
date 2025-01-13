package application

import (
	"GophKeeper/internal/client/models"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rivo/tview"
)

// addRegForm форма регистрации
func (app *Application) addRegForm() *tview.Form {
	data := models.InputDataUser{}
	var password2 string

	app.regForm.AddInputField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "login",
		},
	}), "", 20, nil, func(val string) {
		data.Login = val
	})
	app.regForm.AddPasswordField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "password",
		},
	}), "", 20, 0, func(val string) {
		data.Password = val
	})
	app.regForm.AddPasswordField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "passwordRetry",
		},
	}), "", 20, 0, func(val string) {
		password2 = val
	})
	app.regForm.AddButton(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "save",
		},
	}), func() {
		if err := validateRegForm(data.Login, data.Password, password2); err != nil {
			app.regModal.SetText(err.Error())
		} else {
			if err := app.transport.Register(app.ctx, data); err != nil {
				app.regModal.SetText(err.Error())
			} else {
				app.regModal.SetText(app.localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: "success",
					},
				}))
			}
		}
		app.pages.SwitchToPage(RegModalLink)
	})
	app.regForm.AddButton(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "quit",
		},
	}), func() {
		app.pages.SwitchToPage(MenuLink)
	})
	return app.regForm
}
