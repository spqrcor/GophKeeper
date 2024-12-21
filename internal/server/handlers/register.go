package handlers

import (
	"GophKeeper/internal/server/storage"
	"GophKeeper/internal/server/utils"
	"errors"
	"net/http"
)

// RegisterHandler обработчик роута: POST /api/user/register
func RegisterHandler(s storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var input storage.InputDataUser
		if err := utils.FromPostJSON(req, &input); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err := s.Register(req.Context(), input)
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
