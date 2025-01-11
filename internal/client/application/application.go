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

	app.initModals()
	app.newItemsMenu = app.createNewItemMenu()
	app.menu = app.createMainMenu()

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
func (app *Application) syncData() {
	items, err := app.transport.GetItems(app.ctx)
	if err != nil {
		return
	}
	app.data = items
	app.refreshItemsList()
}

// refreshItemsList выставление списка
func (app *Application) refreshItemsList() {
	app.itemsList.Clear()
	for index, item := range app.data {
		app.itemsList.AddItem(formatTitle(item), " ", rune(49+index), nil)
	}
}

// setConcatText отображение подробных данных
func (app *Application) setConcatText(item *models.ItemData) {
	app.itemText.Clear()
	app.itemText.SetText(formatFull(*item, app.config, app.transport.GetData()))
}

// deleteListItem удаление записи
func (app *Application) deleteListItem(i int) {
	if err := app.transport.RemoveItem(app.ctx, app.data[i].Id); err != nil {
		app.removeModal.SetText(err.Error())
		app.pages.SwitchToPage(RemoveModalLink)
	} else {
		app.data = append(app.data[:i], app.data[i+1:]...)
		app.itemText.Clear()
		app.refreshItemsList()
	}
}

// Start запуск приложения
func (app *Application) Start() {
	transportData := app.transport.GetData()
	if len(transportData.Token) > 0 {
		app.syncData()
	}

	app.itemsList.SetSelectedFunc(func(index int, name string, second_name string, shortcut rune) {
		app.setConcatText(&app.data[index])
	})

	app.flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(app.itemsList, 0, 1, true).
			AddItem(app.itemText, 0, 4, false), 0, 6, false).
		AddItem(app.text, 0, 1, false)

	app.flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			app.tapp.Stop()
		} else if event.Rune() == 98 {
			app.pages.SwitchToPage(MenuLink)
		} else if event.Rune() == 114 {
			app.deleteListItem(app.itemsList.GetCurrentItem())
		}
		return event
	})
	app.addPages()

	if err := app.tapp.SetRoot(app.pages, true).EnableMouse(true).Run(); err != nil {
		app.logger.Fatal(err.Error())
	}
}
