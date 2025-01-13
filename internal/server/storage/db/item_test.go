package db

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/logger"
	"GophKeeper/internal/server/utils"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestCreateItemInfoDB(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	dbres, _, _ := sqlmock.New()
	obj := CreateItemInfoDB(conf, loggerRes, dbres)
	assert.Equal(t, reflect.TypeOf(obj).String() == "db.ItemInfoDB", true)
}

func TestItemInfoDB_GetItem(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := ItemInfoDB{
		config: conf,
		logger: loggerRes,
		DB:     db,
	}
	mock.ExpectExec(getItemQuery).WithArgs("672124b6-9894-11e5-be38-001d42e813fe", utils.CreateKeyFromPin("1234", conf.Salt), "672124b6-9894-11e5-be38-001d42e813fe")
	_, _, _ = d.GetItem(context.Background(), "672124b6-9894-11e5-be38-001d42e813fe", "672124b6-9894-11e5-be38-001d42e813fe", "1234")
}
