package application

import (
	"GophKeeper/internal/client/models"
	"fmt"
)

const (
	minPasswordLength = 6
	maxPasswordLength = 72
	minLoginLength    = 3
	maxLoginLength    = 255
)

var (
	ErrNotEqual     = fmt.Errorf("значения не равны")
	ErrPinMinLength = fmt.Errorf("минимальная длина PIN - 4 знака")
	ErrValidation   = fmt.Errorf("validation error")
	ErrCardFormat   = fmt.Errorf("ошибка валидации номера карты")
)

// validatePinForm валидация формы с пином
func validatePinForm(pin string, pin2 string) error {
	if pin != pin2 {
		return ErrNotEqual
	}
	if len(pin) < 4 {
		return ErrPinMinLength
	}
	return nil
}

// validateNewItem валидация новой записи
func validateNewItem(data models.ItemData) error {
	if data.Type == "CARD" && !luhnAlgorithm(data.CardNum) {
		return ErrCardFormat
	}
	return nil
}

// luhnAlgorithm валидация по алгоритму Луна
func luhnAlgorithm(cardNumber string) bool {
	total := 0
	isSecondDigit := false
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')
		if isSecondDigit {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		total += digit
		isSecondDigit = !isSecondDigit
	}
	return total%10 == 0
}

// validateRegForm валидация формы регистрации
func validateRegForm(login string, password string, password2 string) error {
	if password != password2 {
		return ErrNotEqual
	}
	if (len(login) < minLoginLength) || (len(login) > maxLoginLength) {
		return fmt.Errorf("%w: ошибка при заполнении login, корректная длина от 3 до 255", ErrValidation)
	}
	if (len(password) < minPasswordLength) || (len(password) > maxPasswordLength) {
		return fmt.Errorf("%w: ошибка при заполнении password, корректная длина от 6 до 72", ErrValidation)
	}
	return nil
}
