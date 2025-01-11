package application

import (
	"GophKeeper/internal/client/models"
	"github.com/rivo/tview"
)

// addNewItemForm форма для добавления
func (app *Application) addNewItemForm(itemType string) *tview.Form {
	var filepath string
	data := models.ItemData{
		Type: itemType,
	}

	if itemType == "TEXT" {
		app.newItemForm.AddInputField("Text", "", 40, nil, func(val string) {
			data.Text = val
		})
	}

	if itemType == "AUTH" {
		app.newItemForm.AddInputField("Login", "", 20, nil, func(val string) {
			data.Login = val
		})
		app.newItemForm.AddPasswordField("Password", "", 20, 0, func(val string) {
			data.Password = val
		})
	}

	if itemType == "FILE" {
		app.newItemForm.AddInputField("File", "", 40, nil, func(val string) {
			filepath = val
		})
	}

	if itemType == "CARD" {
		app.newItemForm.AddInputField("Card Number", "", 20, nil, func(val string) {
			data.CardNum = val
		})
		app.newItemForm.AddInputField("Valid Date", "", 10, nil, func(val string) {
			data.CardValid = val
		})
		app.newItemForm.AddInputField("CVV", "", 10, nil, func(val string) {
			data.CardPin = val
		})
	}

	if itemType == "FILE" {
		app.newItemForm.AddButton("save", func() {
			if fileData, err := app.transport.AddItemFile(app.ctx, filepath); err != nil {
				app.newItemModal.SetText(err.Error())
			} else {
				app.data = append(app.data, fileData)
				app.refreshItemsList()
				app.newItemModal.SetText("Успешное добавление")
			}
			app.pages.SwitchToPage(NewItemModalLink)

		})
	} else {
		app.newItemForm.AddButton("save", func() {
			if uuid, err := app.transport.AddItem(app.ctx, data); err != nil {
				app.newItemModal.SetText(err.Error())
			} else {
				data.Id = uuid
				app.data = append(app.data, data)
				app.refreshItemsList()
				app.newItemModal.SetText("Успешное добавление")
			}
			app.pages.SwitchToPage(NewItemModalLink)

		})
	}

	app.newItemForm.AddButton("quit", func() {
		app.pages.SwitchToPage(NewItemsMenuLink)
	})
	return app.newItemForm
}
