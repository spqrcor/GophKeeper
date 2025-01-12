package rest

import (
	"GophKeeper/internal/client/models"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

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
