package application

import (
	"github.com/rivo/tview"
)

const (
	MenuLink         = "Menu"
	PinFormLink      = "Pin Form"
	RegFormLink      = "Reg Form"
	LoginFormLink    = "Login Form"
	NewItemsMenuLink = "New Items Menu"
	PinModalLink     = "Pin Modal"
	RegModalLink     = "Reg Modal"
	LoginModalLink   = "Login Modal"
	RemoveModalLink  = "Remove Modal"
	NewItemModalLink = "New Item Modal"
	NewItemFormLink  = "New Item Form"
	ItemsLink        = "Items"
)

// menuItem Пункт в меню
type menuItem struct {
	MainText      string
	SecondaryText string
	ShortCut      rune
	Selected      func()
}

// pageItem Страница
type pageItem struct {
	Name    string
	Item    tview.Primitive
	Resize  bool
	Visible bool
}

// modalItem модальное окно
type modalItem struct {
	modal *tview.Modal
	link  string
}

// initModals инициализация модальных окон
func (app *Application) initModals() {
	modals := []modalItem{
		{
			app.pinModal,
			PinFormLink,
		},
		{
			app.regModal,
			RegFormLink,
		},
		{
			app.loginModal,
			LoginFormLink,
		},
		{
			app.newItemModal,
			NewItemFormLink,
		},
		{
			app.removeModal,
			ItemsLink,
		},
	}

	for _, item := range modals {
		item.modal = tview.NewModal().
			AddButtons([]string{"Ok"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.pages.SwitchToPage(item.link)
		})
	}
}

// addPages Добавление страниц
func (app *Application) addPages() {
	pages := []pageItem{
		{
			Name:    MenuLink,
			Item:    app.menu,
			Resize:  true,
			Visible: true,
		},
		{
			Name:    PinFormLink,
			Item:    app.pinForm,
			Resize:  true,
			Visible: false,
		},
		{
			Name:    RegFormLink,
			Item:    app.regForm,
			Resize:  true,
			Visible: false,
		},
		{
			Name:    LoginFormLink,
			Item:    app.loginForm,
			Resize:  true,
			Visible: false,
		},
		{
			Name:    NewItemsMenuLink,
			Item:    app.newItemsMenu,
			Resize:  true,
			Visible: false,
		},
		{
			Name:    PinModalLink,
			Item:    app.pinModal,
			Resize:  false,
			Visible: false,
		},
		{
			Name:    RegModalLink,
			Item:    app.regModal,
			Resize:  false,
			Visible: false,
		},
		{
			Name:    LoginModalLink,
			Item:    app.loginModal,
			Resize:  true,
			Visible: false,
		},
		{
			Name:    RemoveModalLink,
			Item:    app.removeModal,
			Resize:  false,
			Visible: false,
		},
		{
			Name:    NewItemModalLink,
			Item:    app.newItemModal,
			Resize:  false,
			Visible: false,
		},
		{
			Name:    NewItemFormLink,
			Item:    app.newItemForm,
			Resize:  true,
			Visible: false,
		},
		{
			Name:    ItemsLink,
			Item:    app.flex,
			Resize:  true,
			Visible: false,
		},
	}

	for _, item := range pages {
		app.pages.AddPage(item.Name, item.Item, item.Resize, item.Visible)
	}
}

// createMainMenu создание главного меню
func (app *Application) createMainMenu() *tview.List {
	menu := tview.NewList()
	mainMenu := []menuItem{
		{
			MainText:      "Set PIN",
			SecondaryText: "Установка PIN",
			ShortCut:      'p',
			Selected: func() {
				app.pinForm.Clear(true)
				app.addPinForm()
				app.pages.SwitchToPage(PinFormLink)
			},
		},
		{
			MainText:      "Registration",
			SecondaryText: "Форма регистрации",
			ShortCut:      'r',
			Selected: func() {
				app.regForm.Clear(true)
				app.addRegForm()
				app.pages.SwitchToPage(RegFormLink)
			},
		},
		{
			MainText:      "Login form",
			SecondaryText: "Форма авторизации",
			ShortCut:      'l',
			Selected: func() {
				app.loginForm.Clear(true)
				app.addLoginForm()
				app.pages.SwitchToPage(LoginFormLink)
			},
		},
		{
			MainText:      "Items",
			SecondaryText: "Список записей",
			ShortCut:      'i',
			Selected: func() {
				app.pages.SwitchToPage(ItemsLink)
			},
		},
		{
			MainText:      "New item",
			SecondaryText: "Добавление записи",
			ShortCut:      'n',
			Selected: func() {
				app.pages.SwitchToPage(NewItemsMenuLink)
			},
		},
		{
			MainText:      "Quit",
			SecondaryText: "Закрыть приложение",
			ShortCut:      'q',
			Selected: func() {
				app.tapp.Stop()
			},
		},
	}

	for _, item := range mainMenu {
		menu.AddItem(item.MainText, item.SecondaryText, item.ShortCut, item.Selected)
	}
	return menu
}

// createNewItemMenu создание меню с выбором типа записи
func (app *Application) createNewItemMenu() *tview.List {
	menu := tview.NewList()
	newItemMenu := []menuItem{
		{
			MainText:      "New Item Text",
			SecondaryText: "Добавление текста",
			ShortCut:      't',
			Selected: func() {
				app.newItemForm.Clear(true)
				app.addNewItemForm("TEXT")
				app.pages.SwitchToPage(NewItemFormLink)
			},
		},
		{
			MainText:      "New Item Login",
			SecondaryText: "Добавление логина/пароля",
			ShortCut:      'l',
			Selected: func() {
				app.newItemForm.Clear(true)
				app.addNewItemForm("AUTH")
				app.pages.SwitchToPage(NewItemFormLink)
			},
		},
		{
			MainText:      "New Item Card",
			SecondaryText: "Добавление данных по карте",
			ShortCut:      'c',
			Selected: func() {
				app.newItemForm.Clear(true)
				app.addNewItemForm("CARD")
				app.pages.SwitchToPage(NewItemFormLink)
			},
		},
		{
			MainText:      "New Item File",
			SecondaryText: "Добавление файла",
			ShortCut:      'f',
			Selected: func() {
				app.newItemForm.Clear(true)
				app.addNewItemForm("FILE")
				app.pages.SwitchToPage(NewItemFormLink)
			},
		},
		{
			MainText:      "Back",
			SecondaryText: "Вернуться в главное меню",
			ShortCut:      'b',
			Selected: func() {
				app.pages.SwitchToPage(MenuLink)
			},
		},
		{
			MainText:      "Quit",
			SecondaryText: "Закрыть приложение",
			ShortCut:      'q',
			Selected: func() {
				app.tapp.Stop()
			},
		},
	}

	for _, item := range newItemMenu {
		menu.AddItem(item.MainText, item.SecondaryText, item.ShortCut, item.Selected)
	}
	return menu
}
