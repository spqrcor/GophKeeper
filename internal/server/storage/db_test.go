package storage

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/logger"
	"GophKeeper/internal/server/utils"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
	"testing"
)

func TestDBStorage_Register(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := DBStorage{
		config: conf,
		logger: loggerRes,
		DB:     db,
	}
	mock.ExpectExec(addUserQuery).WithArgs("672124b6-9894-11e5-be38-001d42e813fe", "spqr", "123")

	_ = d.Register(context.Background(), InputDataUser{Login: "spqr", Password: "123"})
	_ = d.Register(context.Background(), InputDataUser{Login: "spqr", Password: "123456"})
}

func TestDBStorage_Login(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := DBStorage{
		config: conf,
		logger: loggerRes,
		DB:     db,
	}
	mock.ExpectExec(getUserByLoginQuery).WithArgs("spqr")
	_, _ = d.Login(context.Background(), InputDataUser{Login: "spqr", Password: "123"})
}

func TestDBStorage_ShutDown(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := DBStorage{
		config: conf,
		logger: loggerRes,
		DB:     db,
	}
	_ = d.ShutDown()
}

func TestDBStorage_GetItems(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := DBStorage{
		config: conf,
		logger: loggerRes,
		DB:     db,
	}
	mock.ExpectExec(getAllItemsQuery).WithArgs("672124b6-9894-11e5-be38-001d42e813fe", utils.CreateKeyFromPin("spqr", conf.Salt))
	_, _ = d.GetItems(context.Background(), "672124b6-9894-11e5-be38-001d42e813fe", "spqr")
}

func TestDBStorage_RemoveItem(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := DBStorage{
		config: conf,
		logger: loggerRes,
		DB:     db,
	}
	mock.ExpectExec(removeItemQuery).WithArgs("672124b6-9894-11e5-be38-001d42e813fe", "672124b6-9894-11e5-be38-001d42e813fe")
	_ = d.RemoveItem(context.Background(), "672124b6-9894-11e5-be38-001d42e813fe", "672124b6-9894-11e5-be38-001d42e813fe")
}

func TestDBStorage_GetItem(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := DBStorage{
		config: conf,
		logger: loggerRes,
		DB:     db,
	}
	mock.ExpectExec(getItemQuery).WithArgs("672124b6-9894-11e5-be38-001d42e813fe", utils.CreateKeyFromPin("1234", conf.Salt), "672124b6-9894-11e5-be38-001d42e813fe")
	_, _, _ = d.GetItem(context.Background(), "672124b6-9894-11e5-be38-001d42e813fe", "672124b6-9894-11e5-be38-001d42e813fe", "1234")
}

func TestDBStorage_AddItem(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := DBStorage{
		config: conf,
		logger: loggerRes,
		DB:     db,
	}
	mock.ExpectExec(getItemQuery).WithArgs("672124b6-9894-11e5-be38-001d42e813fe", "data", utils.CreateKeyFromPin("1234", conf.Salt))
	_, _ = d.AddItem(context.Background(), CommonData{Type: "TEXT", Text: "1234"}, "672124b6-9894-11e5-be38-001d42e813fe", "1234", nil)
}
