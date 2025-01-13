package application

import (
	"GophKeeper/internal/server/config"
	"GophKeeper/internal/server/logger"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"reflect"
	"syscall"
	"testing"
	"time"
)

func TestNewApplication(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	dbRes, _, _ := sqlmock.New()

	server := NewApplication(
		WithLogger(loggerRes),
		WithConfig(conf),
		WithDB(dbRes),
		WithTokenAuth(jwtauth.New("HS256", []byte(conf.SecretKey), nil, jwt.WithAcceptableSkew(conf.TokenExp))),
	)
	assert.Equal(t, reflect.TypeOf(server).String() == "*application.Application", true)
}

func TestStart(t *testing.T) {
	conf := config.NewConfig()
	loggerRes, _ := logger.NewLogger(zap.InfoLevel)
	dbRes, _, _ := sqlmock.New()

	server := NewApplication(
		WithLogger(loggerRes),
		WithConfig(conf),
		WithDB(dbRes),
		WithTokenAuth(jwtauth.New("HS256", []byte(conf.SecretKey), nil, jwt.WithAcceptableSkew(conf.TokenExp))),
	)
	go func() {
		time.Sleep(1 * time.Second)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	server.Start()
}
