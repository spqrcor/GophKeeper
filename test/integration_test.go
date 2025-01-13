package tests

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/db"
	"GophKeeper/internal/server/logger"
	"GophKeeper/internal/server/storage"
	"context"
	"crypto/tls"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	login := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 10)
	pin := gofakeit.LetterN(8)

	conf := config.NewConfig()
	log, _ := logger.NewLogger(conf.LogLevel)
	dbRes, err := db.Connect(conf.DatabaseDSN)
	if err != nil {
		log.Fatal(err.Error())
	}

	u := url.URL{
		Scheme:     "https",
		Host:       conf.Addr,
		ForceQuery: true,
	}

	e := httpexpect.WithConfig(httpexpect.Config{
		TestName: t.Name(),
		BaseURL:  u.String(),
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCurlPrinter(t),
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	t.Run("register not post", func(t *testing.T) {
		e.GET("/api/user/register").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("register not json", func(t *testing.T) {
		e.POST("/api/user/register").
			WithText("tes").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("register short password", func(t *testing.T) {
		e.POST("/api/user/register").
			WithJSON(storage.InputDataUser{
				Login:    login,
				Password: "333",
			}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("register short login", func(t *testing.T) {
		e.POST("/api/user/register").
			WithJSON(storage.InputDataUser{
				Login:    "l",
				Password: password,
			}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("success register", func(t *testing.T) {
		e.POST("/api/user/register").
			WithJSON(storage.InputDataUser{
				Login:    login,
				Password: password,
			}).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("user exists", func(t *testing.T) {
		e.POST("/api/user/register").
			WithJSON(storage.InputDataUser{
				Login:    login,
				Password: password,
			}).
			Expect().
			Status(http.StatusConflict)
	})

	t.Run("bad login or password", func(t *testing.T) {
		e.POST("/api/user/login").
			WithJSON(storage.InputDataUser{
				Login:    login,
				Password: "333",
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("bad login, empty pin", func(t *testing.T) {
		e.POST("/api/user/login").
			WithJSON(storage.InputDataUser{
				Login:    login,
				Password: password,
			}).
			Expect().
			Status(http.StatusBadRequest)
	})

	var token string
	t.Run("success login", func(t *testing.T) {
		e.POST("/api/user/login").
			WithJSON(storage.InputDataUser{
				Login:    login,
				Password: password,
				Pin:      pin,
			}).
			Expect().
			Status(http.StatusOK).
			Text().Decode(&token)

	})

	t.Run("add item 401 error", func(t *testing.T) {
		e.POST("/api/items").
			WithJSON(storage.CommonData{
				Type: token,
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("add text item format error", func(t *testing.T) {
		e.POST("/api/items").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(storage.CommonData{
				Type: "TEXT",
			}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("add text item success", func(t *testing.T) {
		e.POST("/api/items").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(storage.CommonData{
				Type: "TEXT",
				Text: "test",
			}).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("get items 401", func(t *testing.T) {
		e.GET("/api/items").
			Expect().
			Status(http.StatusUnauthorized)
	})

	var items []storage.CommonData
	t.Run("get items success", func(t *testing.T) {
		e.GET("/api/items").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).JSON().
			Decode(&items)
	})

	if len(items) != 1 || items[0].Text != "test" {
		t.Fail()
	}

	t.Run("add auth item format error", func(t *testing.T) {
		e.POST("/api/items").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(storage.CommonData{
				Type:  "AUTH",
				Login: "test",
			}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("add auth item success", func(t *testing.T) {
		e.POST("/api/items").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(storage.CommonData{
				Type:     "AUTH",
				Login:    "test",
				Password: "test",
			}).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("add card item format error", func(t *testing.T) {
		e.POST("/api/items").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(storage.CommonData{
				Type:  "CARD",
				Login: "test",
			}).
			Expect().
			Status(http.StatusBadRequest)
	})

	var itemID string
	t.Run("add card item success", func(t *testing.T) {
		e.POST("/api/items").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(storage.CommonData{
				Type:      "CARD",
				CardNum:   "0000 0000 0000 0000",
				CardValid: "12/32",
				CardPin:   "000",
			}).
			Expect().
			Status(http.StatusOK).
			Text().Decode(&itemID)
	})

	t.Run("remove item error 401", func(t *testing.T) {
		e.DELETE("/api/items/" + itemID).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("remove item success", func(t *testing.T) {
		e.DELETE("/api/items/"+itemID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK)
	})

	var itemWithFileID string
	t.Run("add file item success", func(t *testing.T) {
		e.POST("/api/items/file").
			WithHeader("Authorization", "Bearer "+token).
			WithMultipart().
			WithFile("file", "../data/test/test1.txt").
			Expect().
			Status(http.StatusOK).
			Text().Decode(&itemWithFileID)
	})

	t.Run("get item file error", func(t *testing.T) {
		e.GET("/api/items/file/" + itemWithFileID + "/token/777777777777").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("get item file success", func(t *testing.T) {
		e.GET("/api/items/file/" + itemWithFileID + "/token/" + token).
			Expect().
			Status(http.StatusOK)
	})

	childCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_ = dbRes.QueryRowContext(childCtx, "DELETE FROM users WHERE login = $1", login)

}
