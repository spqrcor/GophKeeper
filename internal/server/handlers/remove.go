package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

// RemoveItem интерфейс удаление записи
type RemoveItem interface {
	RemoveItem(ctx context.Context, userID string, itemId string) error
}

// RemoveItemHandler обработчик роута: DELETE /api/items/{id}
func RemoveItemHandler(r RemoveItem) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		_, claims, err := jwtauth.FromContext(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = r.RemoveItem(context.Background(), claims["UserID"].(string), chi.URLParam(req, "id")); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}
