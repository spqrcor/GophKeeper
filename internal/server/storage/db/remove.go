package db

import (
	"GophKeeper/internal/server/config"
	"context"
	"database/sql"
	"time"
)

const removeItemQuery = "DELETE FROM user_data WHERE user_id = $1 and id = $2"

// RemoveItemDB тип удаления записи в db
type RemoveItemDB struct {
	config config.Config
	DB     *sql.DB
}

// CreateRemoveItemDB создание RemoveItemDB
func CreateRemoveItemDB(config config.Config, res *sql.DB) RemoveItemDB {
	return RemoveItemDB{
		config: config,
		DB:     res,
	}
}

// RemoveItem удаление записи
func (d RemoveItemDB) RemoveItem(ctx context.Context, userID string, itemId string) error {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*d.config.QueryTimeOut)
	defer cancel()
	_, err := d.DB.ExecContext(childCtx, removeItemQuery, userID, itemId)
	if err != nil {
		return err
	}
	return nil
}
