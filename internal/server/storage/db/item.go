package db

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/storage"
	"GophKeeper/internal/server/utils"
	"context"
	"database/sql"
	"encoding/json"
	"go.uber.org/zap"
	"time"
)

const getItemQuery = "SELECT id, created_at, pgp_sym_decrypt(data,$2,'compress-algo=1, cipher-algo=aes256') as data, pgp_sym_decrypt_bytea(file,$2,'compress-algo=1, cipher-algo=aes256') as file FROM user_data WHERE user_id = $1 and id = $3"

// ItemInfoDB тип информации о записи в db
type ItemInfoDB struct {
	config config.Config
	logger *zap.Logger
	DB     *sql.DB
}

// CreateItemInfoDB создание ItemInfoDB
func CreateItemInfoDB(config config.Config, logger *zap.Logger, res *sql.DB) ItemInfoDB {
	return ItemInfoDB{
		config: config,
		logger: logger,
		DB:     res,
	}
}

// GetItem получение записи по id
func (d ItemInfoDB) GetItem(ctx context.Context, userID string, itemId string, pin string) (storage.CommonData, []byte, error) {
	var item storage.CommonData
	var fileBytes []byte
	childCtx, cancel := context.WithTimeout(ctx, time.Second*d.config.QueryTimeOut)
	defer cancel()

	rows, err := d.DB.QueryContext(childCtx, getItemQuery, userID, utils.CreateKeyFromPin(pin, d.config.Salt), itemId)
	if err != nil {
		return item, nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			d.logger.Error(err.Error())
		}
		if err := rows.Err(); err != nil {
			d.logger.Error(err.Error())
		}
	}()

	for rows.Next() {
		data := ""
		if err = rows.Scan(&item.Id, &item.CreatedAt, &data, &fileBytes); err != nil {
			return item, nil, err
		}
		if err = json.Unmarshal([]byte(data), &item); err != nil {
			return item, nil, err
		}
	}
	return item, fileBytes, nil
}
