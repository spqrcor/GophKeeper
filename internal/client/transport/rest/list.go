package rest

import (
	"GophKeeper/internal/client/models"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

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
