package rest

import (
	"GophKeeper/internal/client/models"
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

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
