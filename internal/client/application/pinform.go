package application

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rivo/tview"
)

// addPinForm форма для выставление пина
func (app *Application) addPinForm() *tview.Form {

	var pin, pin2 string

	app.pinForm.AddPasswordField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "pin",
		},
	}), "", 20, 0, func(val string) {
		pin = val
	})
	app.pinForm.AddPasswordField(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "pinRetry",
		},
	}), "", 20, 0, func(val string) {
		pin2 = val
	})
	app.pinForm.AddButton(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "save",
		},
	}), func() {
		if err := validatePinForm(pin, pin2); err != nil {
			app.pinModal.SetText(err.Error())
		} else {
			transportData := app.transport.GetData()
			transportData.Pin = pin
			if err := app.transport.SetData(); err != nil {
				app.pinModal.SetText(err.Error())
			} else {
				app.pinModal.SetText(app.localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: "success",
					},
				}))
			}
		}
		app.pages.SwitchToPage(PinModalLink)
	})
	app.pinForm.AddButton(app.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "quit",
		},
	}), func() {
		app.pages.SwitchToPage(MenuLink)
	})
	return app.pinForm
}
