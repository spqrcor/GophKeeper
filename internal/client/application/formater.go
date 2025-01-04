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

// formatFull полное отображение записи
func formatFull(data models.ItemData, config config.Config, transportData *transport.Data) string {
	text := "Тип: " + getTypeDescription(data.Type) + "\n"
	text += "Добавлен: " + formatDescription(data) + "\n"

	if data.Type == "TEXT" {
		text += "Текст: " + data.Text
	}
	if data.Type == "CARD" {
		text += "Номер: " + data.CardNum + "\n"
		text += "CVV: " + data.CardPin + "\n"
		text += "Срок действия: " + data.CardValid + "\n"
		text += "Плательщик: " + data.CardPayer
	}
	if data.Type == "FILE" {
		text += "Название: " + data.FileName + "\n"
		text += "Ссылка: " + config.Api + "/api/items/file/" + data.Id + "/token/" + transportData.Token
	}
	if data.Type == "AUTH" {
		text += "Логин: " + data.Login + "\n"
		text += "Пароль: " + data.Password
	}
	return text
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
