package db

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/storage"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCreateDBLoginUser(t *testing.T) {
	dbres, _, _ := sqlmock.New()
	obj := CreateLoginUserDB(config.NewConfig(), dbres)
	assert.Equal(t, reflect.TypeOf(obj).String() == "db.LoginUserDB", true)
}

func TestLoginUserDB_Login(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := LoginUserDB{
		config: config.NewConfig(),
		DB:     db,
	}
	mock.ExpectExec(getUserByLoginQuery).WithArgs("spqr")
	_, _ = d.Login(context.Background(), storage.InputDataUser{Login: "spqr", Password: "123"})
}
