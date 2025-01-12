package rest

import (
	"GophKeeper/internal/client/models"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

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
