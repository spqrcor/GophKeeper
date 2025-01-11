// Package storage
package storage

import (
	"fmt"
	"slices"
)

const (
	minPasswordLength = 6
	maxPasswordLength = 72
	minLoginLength    = 3
	maxLoginLength    = 255
)

var itemTypes = []string{"CARD", "FILE", "AUTH", "TEXT"}

var ErrUnknownType = fmt.Errorf("error unknown type")
var ErrRequired = fmt.Errorf("error required field")
var ErrLogin = fmt.Errorf("login error")
var ErrLoginExists = fmt.Errorf("login exists")
var ErrValidation = fmt.Errorf("validation error")
var ErrGeneratePassword = fmt.Errorf("generate password error")

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
