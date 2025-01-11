package application

import (
	"GophKeeper/internal/client/models"
	"github.com/rivo/tview"
)

// addRegForm форма регистрации
func (app *Application) addRegForm() *tview.Form {
	data := models.InputDataUser{}
	var password2 string

	app.regForm.AddInputField("Login", "", 20, nil, func(val string) {
		data.Login = val
	})
	app.regForm.AddPasswordField("Password", "", 20, 0, func(val string) {
		data.Password = val
	})
	app.regForm.AddPasswordField("Password retry", "", 20, 0, func(val string) {
		password2 = val
	})
	app.regForm.AddButton("save", func() {
		if err := validateRegForm(data.Login, data.Password, password2); err != nil {
			app.regModal.SetText(err.Error())
		} else {
			if err := app.transport.Register(app.ctx, data); err != nil {
				app.regModal.SetText(err.Error())
			} else {
				app.regModal.SetText("Успешная регистрация")
			}
		}
		app.pages.SwitchToPage(RegModalLink)
	})
	app.regForm.AddButton("quit", func() {
		app.pages.SwitchToPage(MenuLink)
	})
	return app.regForm
}
