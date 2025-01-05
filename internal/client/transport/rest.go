package transport

import (
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/crypt"
	"GophKeeper/internal/client/models"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var ErrHttpNotOk = fmt.Errorf("http response not ok")

type RestTransport struct {
	config config.Config
	logger *zap.Logger
	Data   *Data
}

// CreateRestTransport создание db хранилища, config - конфиг, logger - логгер
func CreateRestTransport(config config.Config, logger *zap.Logger, data *Data) Transport {
	return RestTransport{
		config: config,
		logger: logger,
		Data:   data,
	}
}

// Register регистрация
func (r RestTransport) Register(ctx context.Context, input models.InputDataUser) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	childCtx, cancel := context.WithTimeout(ctx, r.config.RequestTimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(childCtx, http.MethodPost, r.config.Api+"/api/user/register", bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := getClient().Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return ErrHttpNotOk
	}
	return nil
}

// Login авторизация
func (r RestTransport) Login(ctx context.Context, input models.InputDataUser) (string, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	childCtx, cancel := context.WithTimeout(ctx, r.config.RequestTimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(childCtx, http.MethodPost, r.config.Api+"/api/user/login", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := getClient().Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", ErrHttpNotOk
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "nil", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return string(bodyBytes), nil
}

// AddItem добавление записи
func (r RestTransport) AddItem(ctx context.Context, item models.ItemData) (string, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	childCtx, cancel := context.WithTimeout(ctx, r.config.RequestTimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(childCtx, http.MethodPost, r.config.Api+"/api/items", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+r.Data.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := getClient().Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", ErrHttpNotOk
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "nil", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return string(bodyBytes), nil
}

// AddItemFile добавление файла записи
func (r RestTransport) AddItemFile(ctx context.Context, filePath string) (models.ItemData, error) {
	data := models.ItemData{
		Type:     "FILE",
		FileName: filepath.Base(filePath),
	}

	file, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", data.FileName)
	if err != nil {
		return data, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return data, err
	}
	err = writer.Close()
	if err != nil {
		return data, err
	}

	childCtx, cancel := context.WithTimeout(ctx, r.config.RequestTimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(childCtx, http.MethodPost, r.config.Api+"/api/items/file", body)
	if err != nil {
		return data, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+r.Data.Token)

	resp, err := getClient().Do(req)
	if err != nil {
		return data, err
	}
	if resp.StatusCode != http.StatusOK {
		return data, ErrHttpNotOk
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data.Id = string(bodyBytes)

	return data, nil
}

// GetItems получение значений
func (r RestTransport) GetItems(ctx context.Context) ([]models.ItemData, error) {
	var data []models.ItemData
	childCtx, cancel := context.WithTimeout(ctx, r.config.RequestTimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(childCtx, http.MethodGet, r.config.Api+"/api/items", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+r.Data.Token)

	resp, err := getClient().Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrHttpNotOk
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return data, nil
}

// RemoveItem удаление записи
func (r RestTransport) RemoveItem(ctx context.Context, itemId string) error {
	childCtx, cancel := context.WithTimeout(ctx, r.config.RequestTimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(childCtx, http.MethodDelete, r.config.Api+"/api/items/"+itemId, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+r.Data.Token)
	resp, err := getClient().Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return ErrHttpNotOk
	}
	return nil
}

// GetData получение настроек
func (r RestTransport) GetData() *Data {
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
