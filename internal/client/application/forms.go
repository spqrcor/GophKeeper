package application

import (
	"GophKeeper/internal/client/models"
	"github.com/rivo/tview"
)

// addPinForm форма для выставление пина
func (a *Application) addPinForm() *tview.Form {

	var pin, pin2 string

	a.pinForm.AddPasswordField("PIN", "", 20, 0, func(val string) {
		pin = val
	})
	a.pinForm.AddPasswordField("PIN retry", "", 20, 0, func(val string) {
		pin2 = val
	})
	a.pinForm.AddButton("save", func() {
		if err := validatePinForm(pin, pin2); err != nil {
			a.pinModal.SetText(err.Error())
		} else {
			transportData := a.transport.GetData()
			transportData.Pin = pin
			if err := a.transport.SetData(); err != nil {
				a.pinModal.SetText(err.Error())
			} else {
				a.pinModal.SetText("PIN успешно установлен")
			}
		}
		a.pages.SwitchToPage("Pin Modal")
	})
	a.pinForm.AddButton("quit", func() {
		a.pages.SwitchToPage("Menu")
	})
	return a.pinForm
}

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

// addNewItemForm форма для добавления
func (a *Application) addNewItemForm(itemType string) *tview.Form {
	var filepath string
	data := models.ItemData{
		Type: itemType,
	}

	if itemType == "TEXT" {
		a.newItemForm.AddInputField("Text", "", 40, nil, func(val string) {
			data.Text = val
		})
	}

	if itemType == "AUTH" {
		a.newItemForm.AddInputField("Login", "", 20, nil, func(val string) {
			data.Login = val
		})
		a.newItemForm.AddPasswordField("Password", "", 20, 0, func(val string) {
			data.Password = val
		})
	}

	if itemType == "FILE" {
		a.newItemForm.AddInputField("File", "", 40, nil, func(val string) {
			filepath = val
		})
	}

	if itemType == "CARD" {
		a.newItemForm.AddInputField("Card Number", "", 20, nil, func(val string) {
			data.CardNum = val
		})
		a.newItemForm.AddInputField("Valid Date", "", 10, nil, func(val string) {
			data.CardValid = val
		})
		a.newItemForm.AddInputField("CVV", "", 10, nil, func(val string) {
			data.CardPin = val
		})
	}

	if itemType == "FILE" {
		a.newItemForm.AddButton("save", func() {
			if fileData, err := a.transport.AddItemFile(a.ctx, filepath); err != nil {
				a.newItemModal.SetText(err.Error())
			} else {
				a.data = append(a.data, fileData)
				a.refreshItemsList()
				a.newItemModal.SetText("Успешное добавление")
			}
			a.pages.SwitchToPage("New Item Modal")

		})
	} else {
		a.newItemForm.AddButton("save", func() {
			if uuid, err := a.transport.AddItem(a.ctx, data); err != nil {
				a.newItemModal.SetText(err.Error())
			} else {
				data.Id = uuid
				a.data = append(a.data, data)
				a.refreshItemsList()
				a.newItemModal.SetText("Успешное добавление")
			}
			a.pages.SwitchToPage("New Item Modal")

		})
	}

	a.newItemForm.AddButton("quit", func() {
		a.pages.SwitchToPage("New Items Menu")
	})
	return a.newItemForm
}
