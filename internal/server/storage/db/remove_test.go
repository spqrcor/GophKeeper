package db

import (
	"GophKeeper/internal/server/config"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCreateRemoveItemDB(t *testing.T) {
	dbres, _, _ := sqlmock.New()
	obj := CreateRemoveItemDB(config.NewConfig(), dbres)
	assert.Equal(t, reflect.TypeOf(obj).String() == "db.RemoveItemDB", true)
}

func TestRemoveItemDB_RemoveItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := RemoveItemDB{
		config: config.NewConfig(),
		DB:     db,
	}
	mock.ExpectExec(removeItemQuery).WithArgs("672124b6-9894-11e5-be38-001d42e813fe", "672124b6-9894-11e5-be38-001d42e813fe")
	_ = d.RemoveItem(context.Background(), "672124b6-9894-11e5-be38-001d42e813fe", "672124b6-9894-11e5-be38-001d42e813fe")
}
