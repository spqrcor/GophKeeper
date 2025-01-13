package rest

import (
	"GophKeeper/internal/client/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

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
