package handlers

import (
	"GophKeeper/internal/server/storage"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

// GetItemFileHandler обработчик роута: GET /api/items/file/{id}
func GetItemFileHandler(s storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		_, claims, err := jwtauth.FromContext(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		item, fileBytes, err := s.GetItem(context.Background(), claims["UserID"].(string), chi.URLParam(req, "id"), claims["PIN"].(string))
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if fileBytes == nil {
			http.Error(res, "file not found", http.StatusNoContent)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Type", "application/octet-stream")
		res.Header().Set("Content-Disposition", "attachment; filename="+item.FileName)
		_, err = res.Write(fileBytes)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}
