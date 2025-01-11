package handlers

import (
	"GophKeeper/internal/server/storage"
	"GophKeeper/internal/server/utils"
	"context"
	"errors"
	"net/http"
)

// RegisterUser интерфейс регистрации
type RegisterUser interface {
	Register(ctx context.Context, input storage.InputDataUser) error
}

// RegisterHandler обработчик роута: POST /api/user/register
func RegisterHandler(r RegisterUser) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var input storage.InputDataUser
		if err := utils.FromPostJSON(req, &input); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err := r.Register(req.Context(), input)
		if errors.Is(err, storage.ErrValidation) {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, storage.ErrLoginExists) {
			res.WriteHeader(http.StatusConflict)
			return
		}
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}
