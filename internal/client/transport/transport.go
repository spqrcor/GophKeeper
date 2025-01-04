package transport

import (
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/crypt"
	"GophKeeper/internal/client/models"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"os"
)

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

// Data тип данных локального хранения
type Data struct {
	Pin   string `json:"pin"`
	Token string `json:"token"`
}

// NewTransport создание хранилища, config конфиг, logger - логгер
func NewTransport(config config.Config, logger *zap.Logger) Transport {
	err := os.MkdirAll(config.DataPath, 0755)
	if err != nil {
		logger.Fatal(err.Error())
	}
	fileData, err := os.ReadFile(config.DataPath + "/" + config.TransportFile)
	if err != nil {
		file, err := os.OpenFile(config.DataPath+"/"+config.TransportFile, os.O_CREATE, 0666)
		if err != nil {
			logger.Fatal(err.Error())
		}
		if err := file.Close(); err != nil {
			logger.Fatal(err.Error())
		}
	}

	data := Data{}
	if len(fileData) != 0 {
		decryptData, err := crypt.Decrypt(fileData, config.SecretKey)
		if err == nil {
			_ = json.Unmarshal(decryptData, &data)
		}
	}
	return CreateRestTransport(config, logger, &data)
}
