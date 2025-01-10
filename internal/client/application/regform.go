package application

import (
	"GophKeeper/internal/client/models"
	"github.com/rivo/tview"
)

// addRegForm форма регистрации
func (a *Application) addRegForm() *tview.Form {
	data := models.InputDataUser{}
	var password2 string

	a.regForm.AddInputField("Login", "", 20, nil, func(val string) {
		data.Login = val
	})
	a.regForm.AddPasswordField("Password", "", 20, 0, func(val string) {
		data.Password = val
	})
	a.regForm.AddPasswordField("Password retry", "", 20, 0, func(val string) {
		password2 = val
	})
	a.regForm.AddButton("save", func() {
		if err := validateRegForm(data.Login, data.Password, password2); err != nil {
			a.regModal.SetText(err.Error())
		} else {
			if err := a.transport.Register(a.ctx, data); err != nil {
				a.regModal.SetText(err.Error())
			} else {
				a.regModal.SetText("Успешная регистрация")
			}
		}
		a.pages.SwitchToPage("Reg Modal")
	})
	a.regForm.AddButton("quit", func() {
		a.pages.SwitchToPage("Menu")
	})
	return a.regForm
}
