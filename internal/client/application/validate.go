package application

import "fmt"

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
