package itemview

import (
	"GophKeeper/internal/client/models"
	"fmt"
	"time"
)

func FormatTitle(data models.ItemData) string {
	switch data.Type {
	case "TEXT":
		return "Текст: ****"
	case "CARD":
		return fmt.Sprintf("Карта: %s", data.CardNum)
	case "FILE":
		return fmt.Sprintf("Файл: %s", data.FileName)
	case "AUTH":
		return fmt.Sprintf("Логин/пароль: %s/****", data.Login)
	}
	return "Тип не определен"
}

func FormatDescription(data models.ItemData) string {
	date, _ := time.Parse(time.RFC3339, data.CreatedAt)
	return date.Format("15:04 02.01.2006")
}
