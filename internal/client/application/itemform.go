package application

import (
	"GophKeeper/internal/client/models"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rivo/tview"
)

// addNewItemForm форма для добавления
func (app *Application) addNewItemForm(itemType string) *tview.Form {
	var filepath string
	data := models.ItemData{
		Type: itemType,
	}

	switch data.Type {
	case "TEXT":
		app.newItemForm.AddInputField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "text",
			},
		}), "", 40, nil, func(val string) {
			data.Text = val
		})
	case "CARD":
		app.newItemForm.AddInputField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "cardNumber",
			},
		}), "", 20, nil, func(val string) {
			data.CardNum = val
		})
		app.newItemForm.AddInputField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "validDate",
			},
		}), "", 10, nil, func(val string) {
			data.CardValid = val
		})
		app.newItemForm.AddInputField("CVV", "", 10, nil, func(val string) {
			data.CardPin = val
		})
	case "FILE":
		app.newItemForm.AddInputField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "file",
			},
		}), "", 40, nil, func(val string) {
			filepath = val
		})
	case "AUTH":
		app.newItemForm.AddInputField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "login",
			},
		}), "", 20, nil, func(val string) {
			data.Login = val
		})
		app.newItemForm.AddPasswordField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "password",
			},
		}), "", 20, 0, func(val string) {
			data.Password = val
		})
	}

	if itemType == "FILE" {
		app.newItemForm.AddButton(app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "save",
			},
		}), func() {
			if fileData, err := app.transport.AddItemFile(app.ctx, filepath); err != nil {
				app.newItemModal.SetText(err.Error())
			} else {
				app.data = append(app.data, fileData)
				app.refreshItemsList()
				app.newItemModal.SetText(app.localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: "success",
					},
				}))
			}
			app.pages.SwitchToPage(NewItemModalLink)

		})
	} else {
		app.newItemForm.AddButton(app.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "save",
			},
		}), func() {
			if err := validateNewItem(data); err != nil {
				app.newItemModal.SetText(err.Error())
			} else {
				if uuid, err := app.transport.AddItem(app.ctx, data); err != nil {
					app.newItemModal.SetText(err.Error())
				} else {
					data.Id = uuid
					app.data = append(app.data, data)
					app.refreshItemsList()
					app.newItemModal.SetText(app.localizer.MustLocalize(&i18n.LocalizeConfig{
						DefaultMessage: &i18n.Message{
							ID: "success",
						},
					}))
				}
			}
			app.pages.SwitchToPage(NewItemModalLink)

		})
	}

	app.newItemForm.AddButton(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "quit",
		},
	}), func() {
		app.pages.SwitchToPage(NewItemsMenuLink)
	})
	return app.newItemForm
}
