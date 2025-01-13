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

func TestCreateRegisterUserDB(t *testing.T) {
	dbres, _, _ := sqlmock.New()
	obj := CreateRegisterUserDB(config.NewConfig(), dbres)
	assert.Equal(t, reflect.TypeOf(obj).String() == "db.RegisterUserDB", true)
}

func TestRegisterUserDB_Register(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := RegisterUserDB{
		config: config.NewConfig(),
		DB:     db,
	}
	mock.ExpectExec(addUserQuery).WithArgs("672124b6-9894-11e5-be38-001d42e813fe", "spqr", "123")

	_ = d.Register(context.Background(), storage.InputDataUser{Login: "spqr", Password: "123"})
	_ = d.Register(context.Background(), storage.InputDataUser{Login: "spqr", Password: "123456"})
}
