package handlers

import (
	"GophKeeper/internal/server/storage"
	"GophKeeper/internal/server/utils"
	"context"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

// LoginUser интерфейс авторизации
type LoginUser interface {
	Login(ctx context.Context, input storage.InputDataUser) (string, error)
}

// LoginHandler обработчик роута: POST /api/user/login
func LoginHandler(l LoginUser, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var input storage.InputDataUser
		if err := utils.FromPostJSON(req, &input); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		UserID, err := l.Login(req.Context(), input)
		if errors.Is(err, storage.ErrLogin) {
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if input.Pin == "" {
			http.Error(res, "PIN необходимо выставить", http.StatusBadRequest)
			return
		}

		_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"UserID": UserID, "PIN": input.Pin})
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write([]byte(tokenString))
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}
	}
}
