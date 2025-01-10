package application

import (
	"GophKeeper/internal/client/models"
	"github.com/rivo/tview"
)

// addLoginForm форма авторзизации
func (a *Application) addLoginForm() *tview.Form {
	data := models.InputDataUser{}

	a.loginForm.AddInputField("Login", "", 20, nil, func(val string) {
		data.Login = val
	})
	a.loginForm.AddPasswordField("Password", "", 20, 0, func(val string) {
		data.Password = val
	})
	a.loginForm.AddButton("save", func() {
		transportData := a.transport.GetData()
		data.Pin = transportData.Pin
		if token, err := a.transport.Login(a.ctx, data); err != nil {
			a.loginModal.SetText(err.Error())
		} else {
			transportData.Token = token
			if err := a.transport.SetData(); err != nil {
				a.loginModal.SetText(err.Error())
			} else {
				a.syncData()
				a.loginModal.SetText("Успешная авторизация")
			}
		}
		a.pages.SwitchToPage("Login Modal")
	})
	a.loginForm.AddButton("quit", func() {
		a.pages.SwitchToPage("Menu")
	})
	return a.loginForm
}
