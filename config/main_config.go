package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/antonioo83/notifier-server/internal/utils"
	"github.com/caarlos0/env/v6"
	"os"
	"time"
)

type Config struct {
	ServerAddress     string `env:"SERVER_ADDRESS"`
	BaseURL           string `env:"BASE_URL"`
	DatabaseDsn       string `env:"DATABASE_DSN"`
	FilepathToDBDump  string
	RequestTimeoutSec int64
	EnableHTTPS       bool   `env:"ENABLE_HTTPS" json:"enable_https,omitempty"` // Enable HTTPS connection.
	ConfigFilePath    string `env:"CONFIG" json:"config_file_path,omitempty"`   // Filename of the server configurations.
	ServerType        string `env:"SERVER_TYPE" json:"server_type,omitempty"`
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

const (
	HTTPServer = "http"
	GRPCServer = "grpc"
)

func GetConfigSettings(configFromFile *Config) (Config, error) {
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
	if configFromFile != nil {
		if cfg.ServerAddress == "" {
			cfg.ServerAddress = configFromFile.ServerAddress
		}
		if cfg.BaseURL == "" {
			cfg.BaseURL = configFromFile.BaseURL
		}
		if cfg.DatabaseDsn == "" {
			cfg.DatabaseDsn = configFromFile.DatabaseDsn
		}
		if !cfg.EnableHTTPS {
			cfg.EnableHTTPS = configFromFile.EnableHTTPS
		}
		if cfg.ServerType == "" {
			cfg.ServerType = configFromFile.ServerType
		}
	}

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

// LoadConfigFile this method read a server configurations from a file in the json format.
func LoadConfigFile(configFilePath string) (*Config, error) {
	var configFromFile Config

	file, err := os.OpenFile(configFilePath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, fmt.Errorf("unable to open a configuration file: %w", err)
	}
	defer utils.ResourceClose(file)

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get statistic info about configuration file: %w", err)
	}
	filesize := info.Size()
	jsonConfig := make([]byte, filesize)

	_, err = file.Read(jsonConfig)
	if err != nil {
		return nil, fmt.Errorf("i can't read a configuration file: %w", err)
	}

	err = json.Unmarshal(jsonConfig, &configFromFile)
	if err != nil {
		return nil, fmt.Errorf("i can't parse a configuration json file: %w", err)
	}

	return &configFromFile, nil
}
