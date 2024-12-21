// Package config конфигурация
package config

import (
	"flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

// Config тип для хранение конфига
type Config struct {
	Addr         string        `env:"SERVER_ADDRESS" json:"server_address"`
	LogLevel     zapcore.Level `env:"LOG_LEVEL"`
	DatabaseDSN  string        `env:"DATABASE_DSN" json:"database_dsn"`
	QueryTimeOut time.Duration `env:"QUERY_TIME_OUT"`
	SecretKey    string        `env:"SECRET_KEY"`
	Salt         string        `env:"SALT"`
	TokenExp     time.Duration `env:"TOKEN_EXPIRATION"`
}

// cfg переменная конфига
var cfg = Config{
	Addr:         "localhost:8081",
	LogLevel:     zap.InfoLevel,
	DatabaseDSN:  "postgres://postgres:Sp123456@localhost:5432/gophkeeper?sslmode=disable",
	QueryTimeOut: 3,
	SecretKey:    "KLJ-fo3Fksd3fl!=",
	Salt:         "Sph2b@o_=zx",
	TokenExp:     time.Hour * 3,
}

var once sync.Once

// NewConfig получение конфига
func NewConfig() Config {
	once.Do(func() {

		flag.StringVar(&cfg.Addr, "a", cfg.Addr, "address and port to run server")
		flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database dsn")
		flag.Parse()

		serverAddressEnv, findAddress := os.LookupEnv("SERVER_ADDRESS")
		serverDatabaseDSN, findDatabaseDSN := os.LookupEnv("DATABASE_DSN")
		if findAddress {
			cfg.Addr = serverAddressEnv
		}
		if findDatabaseDSN {
			cfg.DatabaseDSN = serverDatabaseDSN
		}
	})
	return cfg
}
