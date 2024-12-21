// Package storage
package storage

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/db"
	"context"
	"fmt"
	"go.uber.org/zap"
	"slices"
)

var itemTypes = []string{"CARD", "FILE", "AUTH", "TEXT"}

var ErrUnknownType = fmt.Errorf("error unknown type")
var ErrRequired = fmt.Errorf("error required field")

// InputDataUser тип входящих данных пользователя
type InputDataUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Pin      string `json:"pin"`
}

// CommonData обобщенный тип для записи
type CommonData struct {
	Id        string `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Type      string `json:"type,omitempty"`
	Login     string `json:"login,omitempty"`
	Password  string `json:"password,omitempty"`
	Text      string `json:"text,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	CardNum   string `json:"card_num,omitempty"`
	CardPayer string `json:"card_payer,omitempty"`
	CardValid string `json:"card_valid,omitempty"`
	CardPin   string `json:"card_pin,omitempty"`
}

// Storage интерфейс хранилища
type Storage interface {
	Register(ctx context.Context, input InputDataUser) error
	Login(ctx context.Context, input InputDataUser) (string, error)
	AddItem(ctx context.Context, item CommonData, userID string, pin string, fileBytes []byte) (string, error)
	GetItems(ctx context.Context, userID string, pin string) ([]CommonData, error)
	GetItem(ctx context.Context, userID string, itemId string, pin string) (CommonData, []byte, error)
	RemoveItem(ctx context.Context, userID string, itemId string) error
	ShutDown() error
}

// NewStorage создание хранилища, config конфиг, logger - логгер
func NewStorage(config config.Config, logger *zap.Logger) Storage {
	res, err := db.Connect(config.DatabaseDSN)
	if err != nil {
		logger.Fatal(err.Error())
	}
	if err := db.Migrate(res); err != nil {
		logger.Fatal(err.Error())
	}
	return CreateDBStorage(config, logger, res)
}

// ItemValidator простая валидация новой записи
func ItemValidator(item CommonData) error {
	if !slices.Contains(itemTypes, item.Type) {
		return ErrUnknownType
	}
	if item.Type == "TEXT" && item.Text == "" {
		return ErrRequired
	}
	if item.Type == "CARD" && item.CardNum == "" {
		return ErrRequired
	}
	if item.Type == "AUTH" && (item.Login == "" || item.Password == "") {
		return ErrRequired
	}
	if item.Type == "FILE" && item.FileName == "" {
		return ErrRequired
	}
	return nil
}

// UserValidator валидация при регистрации
func UserValidator(input InputDataUser) error {
	if (len(input.Login) < minLoginLength) || (len(input.Login) > maxLoginLength) {
		return fmt.Errorf("%w: ошибка при заполнении login, корректная длина от 3 до 255", ErrValidation)
	}
	if (len(input.Password) < minPasswordLength) || (len(input.Password) > maxPasswordLength) {
		return fmt.Errorf("%w: ошибка при заполнении password, корректная длина от 6 до 72", ErrValidation)
	}
	return nil
}
