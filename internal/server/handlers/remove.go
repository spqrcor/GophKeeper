package handlers

import (
	"GophKeeper/internal/server/storage"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

// RemoveItemHandler обработчик роута: DELETE /api/items/{id}
func RemoveItemHandler(s storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		_, claims, err := jwtauth.FromContext(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = s.RemoveItem(context.Background(), claims["UserID"].(string), chi.URLParam(req, "id")); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}
