package db

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/storage"
	"GophKeeper/internal/server/utils"
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

const (
	addItemQuery         = "INSERT INTO user_data (user_id,data) VALUES ($1, pgp_sym_encrypt($2,$3,'compress-algo=1, cipher-algo=aes256')) RETURNING id"
	addItemQueryWithFile = "INSERT INTO user_data (user_id,data,file) VALUES ($1, pgp_sym_encrypt($2,$3,'compress-algo=1, cipher-algo=aes256'), pgp_sym_encrypt_bytea($4,$3,'compress-algo=1, cipher-algo=aes256')) RETURNING id"
)

// AddItemDB тип удаления записи в db
type AddItemDB struct {
	config config.Config
	DB     *sql.DB
}

// CreateAddItemDB создание AddItemDB
func CreateAddItemDB(config config.Config, res *sql.DB) AddItemDB {
	return AddItemDB{
		config: config,
		DB:     res,
	}
}

// AddItem добавление записи
func (d AddItemDB) AddItem(ctx context.Context, item storage.CommonData, userID string, pin string, fileBytes []byte) (string, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	itemID := ""

	childCtx, cancel := context.WithTimeout(ctx, time.Second*d.config.QueryTimeOut)
	defer cancel()
	if len(fileBytes) == 0 {
		err = d.DB.QueryRowContext(childCtx, addItemQuery, userID, data, utils.CreateKeyFromPin(pin, d.config.Salt)).Scan(&itemID)
	} else {
		err = d.DB.QueryRowContext(childCtx, addItemQueryWithFile, userID, data, utils.CreateKeyFromPin(pin, d.config.Salt), fileBytes).Scan(&itemID)
	}
	if err != nil {
		return "", err
	}
	return itemID, nil
}
