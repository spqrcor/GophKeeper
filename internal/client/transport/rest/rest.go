package rest

import (
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/crypt"
	"GophKeeper/internal/client/transport"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
)

var ErrHttpNotOk = fmt.Errorf("http response not ok")

type RestTransport struct {
	config config.Config
	logger *zap.Logger
	Data   *transport.Data
}

// CreateRestTransport создание db хранилища, config - конфиг, logger - логгер
func CreateRestTransport(config config.Config, logger *zap.Logger) RestTransport {
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

	data := &transport.Data{}
	if len(fileData) != 0 {
		decryptData, err := crypt.Decrypt(fileData, config.SecretKey)
		if err == nil {
			_ = json.Unmarshal(decryptData, &data)
		}
	}
	return RestTransport{
		config: config,
		logger: logger,
		Data:   data,
	}
}

// GetData получение настроек
func (r RestTransport) GetData() *transport.Data {
	return r.Data
}

// SetData Запись новых настроек
func (r RestTransport) SetData() error {
	x, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	cryptData, err := crypt.Encrypt(x, r.config.SecretKey)
	if err != nil {
		return err
	}
	if err := os.Truncate(r.config.DataPath+"/"+r.config.TransportFile, 0); err != nil {
		return err
	}
	if err := os.WriteFile(r.config.DataPath+"/"+r.config.TransportFile, cryptData, 0666); err != nil {
		return err
	}
	return nil
}

// getClient получение http клиента без проверки сертификата
func getClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
