package handlers

import (
	"GophKeeper/internal/server/storage"
	"context"
	"github.com/go-chi/jwtauth/v5"
	"io"
	"net/http"
)

// AddItemFileHandler обработчик роута: POST /api/items/file
func AddItemFileHandler(s storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var input storage.CommonData
		var fileBytes []byte
		file, handler, err := req.FormFile("file")
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		defer func() {
			_ = file.Close()
		}()
		input.FileName = handler.Filename
		input.Type = "FILE"
		fileBytes, err = io.ReadAll(file)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		_, claims, err := jwtauth.FromContext(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		itemID, err := s.AddItem(context.Background(), input, claims["UserID"].(string), claims["PIN"].(string), fileBytes)
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
