// Package utils методы общего назначения
package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// FromPostJSON обработка входящего json
func FromPostJSON(req *http.Request, input any) error {
	if req.Method != http.MethodPost || !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		return errors.New("invalid request")
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(buf.Bytes(), &input); err != nil {
		return err
	}
	return nil
}

// CreateKeyFromPin формирование key из pin
func CreateKeyFromPin(pin string, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pin + salt))
	return hex.EncodeToString(hasher.Sum(nil))
}
