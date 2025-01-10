package application

import (
	"GophKeeper/internal/client/models"
	"github.com/rivo/tview"
)

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
