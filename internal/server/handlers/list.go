package handlers

import (
	"GophKeeper/internal/server/storage"
	"context"
	"encoding/json"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

// ListItem интерфейс списка записей
type ListItem interface {
	GetItems(ctx context.Context, userID string, pin string) ([]storage.CommonData, error)
}

// GetItemsHandler обработчик роута: GET /api/items
func GetItemsHandler(l ListItem) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		_, claims, err := jwtauth.FromContext(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		items, err := l.GetItems(context.Background(), claims["UserID"].(string), claims["PIN"].(string))
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(res)
		if err := enc.Encode(items); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		}
	}
}
