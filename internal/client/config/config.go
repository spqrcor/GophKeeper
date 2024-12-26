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
	Api            string        `env:"API_ADDRESS" json:"api_address"`
	LogLevel       zapcore.Level `env:"LOG_LEVEL"`
	DataPath       string        `env:"DATA_PATH"`
	TransportFile  string        `env:"TRANSPORT_FILE"`
	SecretKey      string        `env:"SECRET_KEY"`
	RequestTimeOut time.Duration `env:"REQUEST_TIME_OUT"`
}

// cfg переменная конфига
var cfg = Config{
	Api:            "https://localhost:8081",
	LogLevel:       zap.InfoLevel,
	DataPath:       "./data",
	TransportFile:  "transport",
	SecretKey:      "KLJ-fo3Fksd3fl!=",
	RequestTimeOut: 30 * time.Second,
}

var once sync.Once

// NewConfig получение конфига
func NewConfig() Config {
	once.Do(func() {

		flag.StringVar(&cfg.Api, "a", cfg.Api, "api address and port")
		flag.Parse()

		apiAddressEnv, findAddress := os.LookupEnv("API_ADDRESS")
		if findAddress {
			cfg.Api = apiAddressEnv
		}
	})
	return cfg
}
