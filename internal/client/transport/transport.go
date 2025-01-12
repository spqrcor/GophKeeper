package transport

import (
	"GophKeeper/internal/client/models"
	"context"
)

// Data тип данных локального хранения
type Data struct {
	Pin   string `json:"pin"`
	Token string `json:"token"`
}

type Transport interface {
	Register(ctx context.Context, input models.InputDataUser) error
	Login(ctx context.Context, input models.InputDataUser) (string, error)
	AddItem(ctx context.Context, item models.ItemData) (string, error)
	AddItemFile(ctx context.Context, filePath string) (models.ItemData, error)
	GetItems(ctx context.Context) ([]models.ItemData, error)
	RemoveItem(ctx context.Context, itemId string) error
	GetData() *Data
	SetData() error
}
