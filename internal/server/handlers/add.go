package handlers

import (
	"GophKeeper/internal/server/storage"
	"GophKeeper/internal/server/utils"
	"context"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

// AddItemHandler обработчик роута: POST /api/items
func AddItemHandler(s storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var input storage.CommonData
		if err := utils.FromPostJSON(req, &input); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := storage.ItemValidator(input); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		_, claims, err := jwtauth.FromContext(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		itemID, err := s.AddItem(context.Background(), input, claims["UserID"].(string), claims["PIN"].(string), nil)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte(itemID))
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}
	}
}
