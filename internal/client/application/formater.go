package application

import (
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/models"
	"GophKeeper/internal/client/transport"
	"fmt"
	"time"
)

// formatTitle краткое отображение записи
func formatTitle(data models.ItemData) string {
	itemType := getTypeDescription(data.Type)
	switch data.Type {
	case "TEXT":
		return fmt.Sprintf("%s: *****", itemType)
	case "CARD":
		return fmt.Sprintf("%s: %s", itemType, data.CardNum)
	case "FILE":
		return fmt.Sprintf("%s: %s", itemType, data.FileName)
	case "AUTH":
		return fmt.Sprintf("%s: %s/****", itemType, data.Login)
	}
	return "Тип не определен"
}

// formatFullText описание текста
func formatFullText(data models.ItemData) string {
	text := "Тип: " + getTypeDescription(data.Type) + "\n"
	text += "Добавлен: " + formatDescription(data) + "\n"
	text += "Текст: " + data.Text
	return text
}

// formatFullCard описание карты
func formatFullCard(data models.ItemData) string {
	text := "Тип: " + getTypeDescription(data.Type) + "\n"
	text += "Добавлен: " + formatDescription(data) + "\n"
	text += "Номер: " + data.CardNum + "\n"
	text += "CVV: " + data.CardPin + "\n"
	text += "Срок действия: " + data.CardValid + "\n"
	text += "Плательщик: " + data.CardPayer
	return text
}

// formatFullFile описание файла
func formatFullFile(data models.ItemData, config config.Config, transportData *transport.Data) string {
	text := "Тип: " + getTypeDescription(data.Type) + "\n"
	text += "Добавлен: " + formatDescription(data) + "\n"
	text += "Название: " + data.FileName + "\n"
	text += "Ссылка: " + config.Api + "/api/items/file/" + data.Id + "/token/" + transportData.Token
	return text
}

// formatFullTextAuth описание авторизации
func formatFullTextAuth(data models.ItemData) string {
	text := "Тип: " + getTypeDescription(data.Type) + "\n"
	text += "Добавлен: " + formatDescription(data) + "\n"
	text += "Логин: " + data.Login + "\n"
	text += "Пароль: " + data.Password
	return text
}

// formatFull полное отображение записи
func formatFull(data models.ItemData, config config.Config, transportData *transport.Data) string {
	switch data.Type {
	case "TEXT":
		return formatFullText(data)
	case "CARD":
		return formatFullCard(data)
	case "FILE":
		return formatFullFile(data, config, transportData)
	case "AUTH":
		return formatFullTextAuth(data)
	}
	return "Тип не определен"
}

// getTypeDescription получение описания типа
func getTypeDescription(itemType string) string {
	switch itemType {
	case "TEXT":
		return "Текст"
	case "CARD":
		return "Карта"
	case "FILE":
		return "Файл"
	case "AUTH":
		return "Логин/пароль"
	}
	return "Тип не определен"
}

// formatDescription форматирование даты
func formatDescription(data models.ItemData) string {
	date, _ := time.Parse(time.RFC3339, data.CreatedAt)
	return date.Format("15:04 02.01.2006")
}
