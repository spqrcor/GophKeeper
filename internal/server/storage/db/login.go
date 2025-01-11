package db

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/storage"
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const getUserByLoginQuery = "SELECT id, password FROM users WHERE login = $1"

// LoginUserDB тип авторизации через db
type LoginUserDB struct {
	config config.Config
	DB     *sql.DB
}

// CreateLoginUserDB создание DBLoginUser
func CreateLoginUserDB(config config.Config, res *sql.DB) LoginUserDB {
	return LoginUserDB{
		config: config,
		DB:     res,
	}
}

// Login авторизация
func (d LoginUserDB) Login(ctx context.Context, input storage.InputDataUser) (string, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*d.config.QueryTimeOut)
	defer cancel()
	row := d.DB.QueryRowContext(childCtx, getUserByLoginQuery, input.Login)

	var userID, password string
	err := row.Scan(&userID, &password)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrLogin
	}
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(input.Password)); err != nil {
		return "", storage.ErrLogin
	}
	return userID, nil
}
