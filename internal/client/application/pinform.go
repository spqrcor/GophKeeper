package application

import (
	"github.com/rivo/tview"
)

// addPinForm форма для выставление пина
func (app *Application) addPinForm() *tview.Form {

	var pin, pin2 string

	app.pinForm.AddPasswordField("PIN", "", 20, 0, func(val string) {
		pin = val
	})
	app.pinForm.AddPasswordField("PIN retry", "", 20, 0, func(val string) {
		pin2 = val
	})
	app.pinForm.AddButton("save", func() {
		if err := validatePinForm(pin, pin2); err != nil {
			app.pinModal.SetText(err.Error())
		} else {
			transportData := app.transport.GetData()
			transportData.Pin = pin
			if err := app.transport.SetData(); err != nil {
				app.pinModal.SetText(err.Error())
			} else {
				app.pinModal.SetText("PIN успешно установлен")
			}
		}
		app.pages.SwitchToPage(PinModalLink)
	})
	app.pinForm.AddButton("quit", func() {
		app.pages.SwitchToPage(MenuLink)
	})
	return app.pinForm
}
