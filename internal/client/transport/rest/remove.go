package rest

import (
	"context"
	"net/http"
)

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
