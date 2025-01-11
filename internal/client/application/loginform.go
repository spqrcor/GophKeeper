package application

import (
	"GophKeeper/internal/client/models"
	"github.com/rivo/tview"
)

// addLoginForm форма авторзизации
func (app *Application) addLoginForm() *tview.Form {
	data := models.InputDataUser{}

	app.loginForm.AddInputField("Login", "", 20, nil, func(val string) {
		data.Login = val
	})
	app.loginForm.AddPasswordField("Password", "", 20, 0, func(val string) {
		data.Password = val
	})
	app.loginForm.AddButton("save", func() {
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
				app.loginModal.SetText("Успешная авторизация")
			}
		}
		app.pages.SwitchToPage(LoginModalLink)
	})
	app.loginForm.AddButton("quit", func() {
		app.pages.SwitchToPage(MenuLink)
	})
	return app.loginForm
}
