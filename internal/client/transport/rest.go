package transport

import (
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/models"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

var ErrHttpNotOk = fmt.Errorf("http response not ok")

type RestTransport struct {
	config config.Config
	logger *zap.Logger
	data   Data
}

// CreateRestTransport создание db хранилища, config - конфиг, logger - логгер
func CreateRestTransport(config config.Config, logger *zap.Logger, data Data) Transport {
	return RestTransport{
		config: config,
		logger: logger,
		data:   data,
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
	defer func() {
		_ = resp.Body.Close()
	}()
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
	req.Header.Set("Authorization", "Bearer "+r.data.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := getClient().Do(req)
	defer func() {
		_ = resp.Body.Close()
	}()
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
	return string(bodyBytes), nil
}

// AddItemFile добавление файла записи
func (r RestTransport) AddItemFile(ctx context.Context, fileBytes []byte) (string, error) {
	return "", nil
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
	req.Header.Set("Authorization", "Bearer "+r.data.Token)

	resp, err := getClient().Do(req)
	defer func() {
		_ = resp.Body.Close()
	}()
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
	return data, nil
}

// GetItemFile получение файла записи
func (r RestTransport) GetItemFile(ctx context.Context, itemId string) ([]byte, error) {
	childCtx, cancel := context.WithTimeout(ctx, r.config.RequestTimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(childCtx, http.MethodGet, r.config.Api+"/api/items/file/"+itemId, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+r.data.Token)
	resp, err := getClient().Do(req)
	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrHttpNotOk
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

// RemoveItem удаление записи
func (r RestTransport) RemoveItem(ctx context.Context, itemId string) error {
	childCtx, cancel := context.WithTimeout(ctx, r.config.RequestTimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(childCtx, http.MethodDelete, r.config.Api+"/api/items/"+itemId, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+r.data.Token)
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
func (r RestTransport) GetData() Data {
	return r.data
}

func getClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
