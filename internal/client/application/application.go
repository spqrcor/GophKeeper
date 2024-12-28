// Package application
package application

import (
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/itemview"
	"GophKeeper/internal/client/models"
	"GophKeeper/internal/client/transport"
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
	"os"
)

// Application тип Application
type Application struct {
	config    config.Config
	logger    *zap.Logger
	transport transport.Transport
	ctx       context.Context
	data      map[string]models.ItemData
}

// NewApplication создание Application, opts набор параметров
func NewApplication(opts ...func(*Application)) *Application {
	app := &Application{
		ctx:  context.Background(),
		data: make(map[string]models.ItemData),
	}
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

var (
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	errorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#d94a2b", Dark: "#d94a2b"}).
				Render
)

type item struct {
	title       string
	description string
	id          string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

type model struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
}

func (a *Application) newModel() model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = &listKeyMap{}
	)

	var items []list.Item
	for _, r := range a.data {
		items = append(items, item{
			title:       itemview.FormatTitle(r),
			description: itemview.FormatDescription(r),
			id:          r.Id,
		})
	}

	delegate := a.newItemDelegate(delegateKeys)
	groceryList := list.New(items, delegate, 0, 0)
	groceryList.Title = "GophKeeper"
	groceryList.Styles.Title = titleStyle

	return model{
		list:         groceryList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

// Start запуск приложения
func (a *Application) Start() {
	transportData := a.transport.GetData()
	if len(transportData.Token) > 0 {
		a.syncData()
	}
	if _, err := tea.NewProgram(a.newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (a *Application) syncData() {
	items, err := a.transport.GetItems(a.ctx)
	if err != nil {
		return
	}
	a.data = map[string]models.ItemData{}
	for _, item := range items {
		a.data[item.Id] = item
	}
}
