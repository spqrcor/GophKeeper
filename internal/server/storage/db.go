package storage

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	minPasswordLength = 6
	maxPasswordLength = 72
	minLoginLength    = 3
	maxLoginLength    = 255
)

var ErrLogin = fmt.Errorf("login error")
var ErrLoginExists = fmt.Errorf("login exists")
var ErrValidation = fmt.Errorf("validation error")
var ErrGeneratePassword = fmt.Errorf("generate password error")

const (
	addUserQuery         = "INSERT INTO users (id, login, password) VALUES ($1, $2, $3) ON CONFLICT(login) DO UPDATE SET login = EXCLUDED.login RETURNING id"
	getUserByLoginQuery  = "SELECT id, password FROM users WHERE login = $1"
	addItemQuery         = "INSERT INTO user_data (user_id,data) VALUES ($1, pgp_sym_encrypt($2,$3,'compress-algo=1, cipher-algo=aes256')) RETURNING id"
	addItemQueryWithFile = "INSERT INTO user_data (user_id,data,file) VALUES ($1, pgp_sym_encrypt($2,$3,'compress-algo=1, cipher-algo=aes256'), pgp_sym_encrypt_bytea($4,$3,'compress-algo=1, cipher-algo=aes256')) RETURNING id"
	getAllItemsQuery     = "SELECT id, created_at, pgp_sym_decrypt(data,$2,'compress-algo=1, cipher-algo=aes256') as data FROM user_data WHERE user_id = $1 ORDER BY created_at DESC"
	getItemQuery         = "SELECT id, created_at, pgp_sym_decrypt(data,$2,'compress-algo=1, cipher-algo=aes256') as data, pgp_sym_decrypt_bytea(file,$2,'compress-algo=1, cipher-algo=aes256') as file FROM user_data WHERE user_id = $1 and id = $3"
	removeItemQuery      = "DELETE FROM user_data WHERE user_id = $1 and id = $2"
)

// DBStorage тип db хранилища
type DBStorage struct {
	config config.Config
	logger *zap.Logger
	DB     *sql.DB
}

// CreateDBStorage создание db хранилища, config - конфиг, logger - логгер
func CreateDBStorage(config config.Config, logger *zap.Logger, res *sql.DB) Storage {
	return DBStorage{
		config: config,
		logger: logger,
		DB:     res,
	}
}

// Register регистрация
func (d DBStorage) Register(ctx context.Context, input InputDataUser) error {
	if err := UserValidator(input); err != nil {
		return err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return errors.Join(ErrGeneratePassword, err)
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
		return ErrLoginExists
	}
	return nil
}

// Login авторизация
func (d DBStorage) Login(ctx context.Context, input InputDataUser) (string, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*d.config.QueryTimeOut)
	defer cancel()
	row := d.DB.QueryRowContext(childCtx, getUserByLoginQuery, input.Login)

	var userID, password string
	err := row.Scan(&userID, &password)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrLogin
	}
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(input.Password)); err != nil {
		return "", ErrLogin
	}
	return userID, nil
}

// ShutDown завершение работы с хранилищем
func (d DBStorage) ShutDown() error {
	return d.DB.Close()
}

// GetItems получение всех записей
func (d DBStorage) GetItems(ctx context.Context, userID string, pin string) ([]CommonData, error) {
	var items []CommonData
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
		item := CommonData{}
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

// RemoveItem удаление записи
func (d DBStorage) RemoveItem(ctx context.Context, userID string, itemId string) error {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*d.config.QueryTimeOut)
	defer cancel()
	_, err := d.DB.ExecContext(childCtx, removeItemQuery, userID, itemId)
	if err != nil {
		return err
	}
	return nil
}

// AddItem добавление записи
func (d DBStorage) AddItem(ctx context.Context, item CommonData, userID string, pin string, fileBytes []byte) (string, error) {
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

// GetItem получение записи по id
func (d DBStorage) GetItem(ctx context.Context, userID string, itemId string, pin string) (CommonData, []byte, error) {
	var item CommonData
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
