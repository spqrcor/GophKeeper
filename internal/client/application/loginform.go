package application

import (
	"GophKeeper/internal/client/models"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rivo/tview"
)

// addLoginForm форма авторзизации
func (app *Application) addLoginForm() *tview.Form {
	data := models.InputDataUser{}

	app.loginForm.AddInputField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "login",
		},
	}), "", 20, nil, func(val string) {
		data.Login = val
	})
	app.loginForm.AddPasswordField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "password",
		},
	}), "", 20, 0, func(val string) {
		data.Password = val
	})
	app.loginForm.AddButton(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "save",
		},
	}), func() {
		transportData := app.transport.GetData()
		data.Pin = transportData.Pin
		if token, err := app.transport.Login(app.ctx, data); err != nil {
			app.loginModal.SetText(err.Error())
		} else {
			transportData.Token = token
			if err := app.transport.SetData(); err != nil {
				app.loginModal.SetText(err.Error())
			} else {
				app.syncData()
				app.loginModal.SetText(app.localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: "success",
					},
				}))
			}
		}
		app.pages.SwitchToPage(LoginModalLink)
	})
	app.loginForm.AddButton(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "quit",
		},
	}), func() {
		app.pages.SwitchToPage(MenuLink)
	})
	return app.loginForm
}
