package application

import (
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
