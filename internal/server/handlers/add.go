package handlers

import (
	"GophKeeper/internal/server/storage"
	"GophKeeper/internal/server/utils"
	"context"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

// AddItem интерфейс данных по записи
type AddItem interface {
	AddItem(ctx context.Context, item storage.CommonData, userID string, pin string, fileBytes []byte) (string, error)
}

// AddItemHandler обработчик роута: POST /api/items
func AddItemHandler(a AddItem) http.HandlerFunc {
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

		itemID, err := a.AddItem(context.Background(), input, claims["UserID"].(string), claims["PIN"].(string), nil)
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
