package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	ServerAddress     string `env:"SERVER_ADDRESS"`
	BaseURL           string `env:"BASE_URL"`
	DatabaseDsn       string `env:"DATABASE_DSN"`
	FilepathToDBDump  string
	RequestTimeoutSec int64
	Auth              Auth
	Callback          Callback
	Sender            Sender
}

type Auth struct {
	AdminAuthToken string
}

type Sender struct {
	ItemCount    int
	MaxAttempts  int
	Timeout      time.Duration
	LoadInterval time.Duration
}

type Callback struct {
	MaxAttempts     uint
	LimitUnitOfTime uint
	CronSpec        string
}

var cfg Config

func GetConfigSettings() (Config, error) {
	const (
		ServerAddress           = ":8080"
		DatabaseDSN             = "postgres://postgres:433370@localhost:5432/notifier_server"
		RequestTimeoutSec       = 60
		AdminAuthToken          = "54d1ba805e2a4891aeac9299b618945e"
		CallbackMaxAttempts     = 3
		CallbackLimitUnitOfTime = 50
		CallbackCronSpec        = "1 * * * * *"
		SenderItemCount         = 10
		SenderMaxAttempts       = 3
		SenderTimeout           = time.Second * 3
		SenderLoadInterval      = time.Second * 7
	)

	err := env.Parse(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("I can't parse config: %w", err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "The address of the local server")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base address of the result short url")
	flag.StringVar(&cfg.DatabaseDsn, "d", cfg.DatabaseDsn, "Database port")
	flag.Parse()
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ServerAddress
	}

	if cfg.DatabaseDsn == "" {
		cfg.DatabaseDsn = DatabaseDSN
	}

	if cfg.RequestTimeoutSec == 0 {
		cfg.RequestTimeoutSec = RequestTimeoutSec
	}

	cfg.Auth.AdminAuthToken = AdminAuthToken

	cfg.Callback.MaxAttempts = CallbackMaxAttempts
	cfg.Callback.LimitUnitOfTime = CallbackLimitUnitOfTime
	cfg.Callback.CronSpec = CallbackCronSpec

	cfg.Sender.ItemCount = SenderItemCount
	cfg.Sender.MaxAttempts = SenderMaxAttempts
	cfg.Sender.Timeout = SenderTimeout
	cfg.Sender.LoadInterval = SenderLoadInterval

	return cfg, nil
}
