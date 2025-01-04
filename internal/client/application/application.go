// Package application
package application

import (
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/models"
	"GophKeeper/internal/client/transport"
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"go.uber.org/zap"
)

// Application тип Application
type Application struct {
	config       config.Config
	logger       *zap.Logger
	transport    transport.Transport
	ctx          context.Context
	data         []models.ItemData
	pages        *tview.Pages
	tapp         *tview.Application
	pinForm      *tview.Form
	regForm      *tview.Form
	loginForm    *tview.Form
	newItemForm  *tview.Form
	menu         *tview.List
	newItemsMenu *tview.List
	pinModal     *tview.Modal
	regModal     *tview.Modal
	loginModal   *tview.Modal
	newItemModal *tview.Modal
	removeModal  *tview.Modal
	itemText     *tview.TextView
	itemsList    *tview.List
	flex         *tview.Flex
	text         *tview.TextView
}

// NewApplication создание Application, opts набор параметров
func NewApplication(opts ...func(*Application)) *Application {
	app := &Application{
		ctx:         context.Background(),
		data:        []models.ItemData{},
		pages:       tview.NewPages(),
		tapp:        tview.NewApplication(),
		pinForm:     tview.NewForm(),
		regForm:     tview.NewForm(),
		loginForm:   tview.NewForm(),
		newItemForm: tview.NewForm(),
		itemText:    tview.NewTextView(),
		itemsList:   tview.NewList().ShowSecondaryText(false),
		flex:        tview.NewFlex(),
		text: tview.NewTextView().
			SetTextColor(tcell.ColorGreen).
			SetText("(r) remove\n(b) back \n(q) to quit"),
	}

	app.pinModal = tview.NewModal().
		AddButtons([]string{"Ok"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.pages.SwitchToPage("Pin Form")
	})
	app.regModal = tview.NewModal().
		AddButtons([]string{"Ok"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.pages.SwitchToPage("Reg Form")
	})
	app.loginModal = tview.NewModal().
		AddButtons([]string{"Ok"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.pages.SwitchToPage("Login Form")
	})
	app.newItemModal = tview.NewModal().
		AddButtons([]string{"Ok"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.pages.SwitchToPage("New Item Form")
	})
	app.removeModal = tview.NewModal().
		AddButtons([]string{"Ok"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.pages.SwitchToPage("Items")
	})

	app.newItemsMenu = tview.NewList().
		AddItem("New Item Text", "Добавление текста", 't', func() {
			app.newItemForm.Clear(true)
			app.addNewItemForm("TEXT")
			app.pages.SwitchToPage("New Item Form")
		}).
		AddItem("New Item Login", "Добавление логина/пароля", 'l', func() {
			app.newItemForm.Clear(true)
			app.addNewItemForm("AUTH")
			app.pages.SwitchToPage("New Item Form")
		}).
		AddItem("New Item Card", "Добавление данных по карте", 'c', func() {
			app.newItemForm.Clear(true)
			app.addNewItemForm("CARD")
			app.pages.SwitchToPage("New Item Form")
		}).
		AddItem("New Item File", "Добавление файла", 'f', func() {
			app.newItemForm.Clear(true)
			app.addNewItemForm("FILE")
			app.pages.SwitchToPage("New Item Form")
		}).
		AddItem("Back", "Вернуться в главное меню", 'b', func() {
			app.pages.SwitchToPage("Menu")
		}).
		AddItem("Quit", "Закрыть приложение", 'q', func() {
			app.tapp.Stop()
		})

	app.menu = tview.NewList().
		AddItem("Set PIN", "Установка PIN", 'p', func() {
			app.pinForm.Clear(true)
			app.addPinForm()
			app.pages.SwitchToPage("Pin Form")
		}).
		AddItem("Registration", "Форма регистрации", 'r', func() {
			app.regForm.Clear(true)
			app.addRegForm()
			app.pages.SwitchToPage("Reg Form")
		}).
		AddItem("Login form", "Форма авторизации", 'l', func() {
			app.loginForm.Clear(true)
			app.addLoginForm()
			app.pages.SwitchToPage("Login Form")
		}).
		AddItem("Items", "Список записей", 'n', func() {
			app.pages.SwitchToPage("Items")
		}).
		AddItem("New item", "Добавление записи", 'n', func() {
			app.pages.SwitchToPage("New Items Menu")
		}).
		AddItem("Quit", "Закрыть приложение", 'q', func() {
			app.tapp.Stop()
		})

	for _, opt := range opts {
		opt(app)
	}
	return app
}

// WithLogger добавление logger
func WithLogger(logger *zap.Logger) func(*Application) {
	return func(a *Application) {
		a.logger = logger
	}
}

// WithConfig добавление config
func WithConfig(config config.Config) func(*Application) {
	return func(a *Application) {
		a.config = config
	}
}

// WithTransport добавление transport
func WithTransport(transport transport.Transport) func(*Application) {
	return func(a *Application) {
		a.transport = transport
	}
}

// syncData получение данных с сервера
func (a *Application) syncData() {
	items, err := a.transport.GetItems(a.ctx)
	if err != nil {
		return
	}
	a.data = items
	a.refreshItemsList()
}

// refreshItemsList выставление списка
func (a *Application) refreshItemsList() {
	a.itemsList.Clear()
	for index, item := range a.data {
		a.itemsList.AddItem(formatTitle(item), " ", rune(49+index), nil)
	}
}

// setConcatText отображение подробных данных
func (a *Application) setConcatText(item *models.ItemData) {
	a.itemText.Clear()
	a.itemText.SetText(formatFull(*item, a.config, a.transport.GetData()))
}

// deleteListItem удаление записи
func (a *Application) deleteListItem(i int) {
	if err := a.transport.RemoveItem(a.ctx, a.data[i].Id); err != nil {
		a.removeModal.SetText(err.Error())
		a.pages.SwitchToPage("Remove Modal")
	} else {
		a.data = append(a.data[:i], a.data[i+1:]...)
		a.itemText.Clear()
		a.refreshItemsList()
	}
}

// Start запуск приложения
func (a *Application) Start() {
	transportData := a.transport.GetData()
	if len(transportData.Token) > 0 {
		a.syncData()
	}

	a.itemsList.SetSelectedFunc(func(index int, name string, second_name string, shortcut rune) {
		a.setConcatText(&a.data[index])
	})

	a.flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(a.itemsList, 0, 1, true).
			AddItem(a.itemText, 0, 4, false), 0, 6, false).
		AddItem(a.text, 0, 1, false)

	a.flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			a.tapp.Stop()
		} else if event.Rune() == 98 {
			a.pages.SwitchToPage("Menu")
		} else if event.Rune() == 114 {
			a.deleteListItem(a.itemsList.GetCurrentItem())
		}
		return event
	})

	a.pages.AddPage("Menu", a.menu, true, true)
	a.pages.AddPage("Pin Form", a.pinForm, true, false)
	a.pages.AddPage("Reg Form", a.regForm, true, false)
	a.pages.AddPage("Login Form", a.loginForm, true, false)
	a.pages.AddPage("New Items Menu", a.newItemsMenu, true, false)
	a.pages.AddPage("Pin Modal", a.pinModal, false, false)
	a.pages.AddPage("Reg Modal", a.regModal, false, false)
	a.pages.AddPage("Login Modal", a.loginModal, false, false)
	a.pages.AddPage("Remove Modal", a.removeModal, false, false)
	a.pages.AddPage("New Item Modal", a.newItemModal, false, false)
	a.pages.AddPage("New Item Form", a.newItemForm, false, false)
	a.pages.AddPage("Items", a.flex, true, false)

	if err := a.tapp.SetRoot(a.pages, true).EnableMouse(true).Run(); err != nil {
		a.logger.Fatal(err.Error())
	}
}
