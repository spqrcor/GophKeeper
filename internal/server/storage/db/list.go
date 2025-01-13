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

const getAllItemsQuery = "SELECT id, created_at, pgp_sym_decrypt(data,$2,'compress-algo=1, cipher-algo=aes256') as data FROM user_data WHERE user_id = $1 ORDER BY created_at DESC"

// ListItemDB тип списка через db
type ListItemDB struct {
	config config.Config
	logger *zap.Logger
	DB     *sql.DB
}

// CreateListItemDB создание ListItemDB
func CreateListItemDB(config config.Config, logger *zap.Logger, res *sql.DB) ListItemDB {
	return ListItemDB{
		config: config,
		logger: logger,
		DB:     res,
	}
}

// GetItems получение всех записей
func (d ListItemDB) GetItems(ctx context.Context, userID string, pin string) ([]storage.CommonData, error) {
	var items []storage.CommonData
	childCtx, cancel := context.WithTimeout(ctx, time.Second*d.config.QueryTimeOut)
	defer cancel()

	rows, err := d.DB.QueryContext(childCtx, getAllItemsQuery, userID, utils.CreateKeyFromPin(pin, d.config.Salt))
	if err != nil {
		return nil, err
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
		item := storage.CommonData{}
		data := ""
		if err = rows.Scan(&item.Id, &item.CreatedAt, &data); err != nil {
			return nil, err
		}
		if err = json.Unmarshal([]byte(data), &item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
