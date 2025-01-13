package db

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/storage"
	"GophKeeper/internal/server/utils"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAddItemDB_AddItem(t *testing.T) {
	conf := config.NewConfig()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	d := AddItemDB{
		config: conf,
		DB:     db,
	}
	mock.ExpectExec(getItemQuery).WithArgs("672124b6-9894-11e5-be38-001d42e813fe", "data", utils.CreateKeyFromPin("1234", conf.Salt))
	_, _ = d.AddItem(context.Background(), storage.CommonData{Type: "TEXT", Text: "1234"}, "672124b6-9894-11e5-be38-001d42e813fe", "1234", nil)

}

func TestCreateAddItemDB(t *testing.T) {
	dbres, _, _ := sqlmock.New()
	obj := CreateAddItemDB(config.NewConfig(), dbres)
	assert.Equal(t, reflect.TypeOf(obj).String() == "db.AddItemDB", true)
}
