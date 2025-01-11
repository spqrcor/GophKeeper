package db

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/storage"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const addUserQuery = "INSERT INTO users (id, login, password) VALUES ($1, $2, $3) ON CONFLICT(login) DO UPDATE SET login = EXCLUDED.login RETURNING id"

// RegisterUserDB тип авторизации через db
type RegisterUserDB struct {
	config config.Config
	DB     *sql.DB
}

// CreateRegisterUserDB создание RegisterUserDB
func CreateRegisterUserDB(config config.Config, res *sql.DB) RegisterUserDB {
	return RegisterUserDB{
		config: config,
		DB:     res,
	}
}

// Register регистрация
func (d RegisterUserDB) Register(ctx context.Context, input storage.InputDataUser) error {
	if err := storage.UserValidator(input); err != nil {
		return err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return errors.Join(storage.ErrGeneratePassword, err)
	}

	baseUserID := ""
	userID := uuid.NewString()
	childCtx, cancel := context.WithTimeout(ctx, time.Second*d.config.QueryTimeOut)
	defer cancel()
	err = d.DB.QueryRowContext(childCtx, addUserQuery, userID, input.Login, string(bytes)).Scan(&baseUserID)
	if err != nil {
		return err
	}
	if userID != baseUserID {
		return storage.ErrLoginExists
	}
	return nil
}
