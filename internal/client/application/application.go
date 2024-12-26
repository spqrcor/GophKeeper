// Package application
package application

import (
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/models"
	"GophKeeper/internal/client/transport"
	"context"
	"go.uber.org/zap"
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

// Start запуск приложения
func (a *Application) Start() {
	transportData := a.transport.GetData()
	if len(transportData.Token) > 0 {
		items, _ := a.transport.GetItems(a.ctx)
		for _, item := range items {
			a.data[item.Id] = item
		}
	}
	//_ = a.transport.RemoveItem(a.ctx, "dafb75a7-6316-4480-bbde-fa337239fa9e")
	//_, _ = a.transport.GetItemFile(a.ctx, "30960cf9-0a69-4168-bfa4-8550e3258871")
	//_, _ = a.transport.AddItem(a.ctx, models.ItemData{Type: "TEXT", Text: "1234"})
	//_ = a.transport.Register(a.ctx, models.InputDataUser{Login: "SPQRCOR", Password: "123456"})
	//_, _ = a.transport.Login(a.ctx, models.InputDataUser{Login: "SPQRCOR", Password: "123456", Pin: transportData.Pin})
	_, _ = a.transport.AddItemFile(a.ctx, []byte("xx"))
}
