package config

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/notifier-server/internal/utils"
	"os"
)

type ClientConfig struct {
	ServerAddress     string `json:"server_address,omitempty"`
	RequestTimeoutSec int64  `json:"request_timeout_sec,omitempty"`
	AdminAuthToken    string `json:"token,omitempty"`
	SendMessageURI    string `json:"send_message_uri,omitempty"`
	StatusGetURI      string `json:"status_get_uri,omitempty"`
}

var clCfg ClientConfig

// GetClientConfigSettings gets configuration settings of the server.
func GetClientConfigSettings(configFromFile *ClientConfig) (ClientConfig, error) {
	const (
		ServerAddress     = "localhost:8080"
		RequestTimeoutSec = 3
		AdminAuthToken    = "54d1ba805e2a4891aeac9299b618945e"
		SendMessageURL    = ServerAddress + "/api/v1/messages"
		GetStatusURL      = ServerAddress + "/api/v1/message"
	)

	clCfg.ServerAddress = configFromFile.ServerAddress
	clCfg.RequestTimeoutSec = configFromFile.RequestTimeoutSec
	clCfg.AdminAuthToken = configFromFile.AdminAuthToken
	clCfg.SendMessageURI = configFromFile.ServerAddress + configFromFile.SendMessageURI
	clCfg.StatusGetURI = configFromFile.ServerAddress + configFromFile.StatusGetURI
	if configFromFile != nil {
		if clCfg.ServerAddress == "" {
			clCfg.ServerAddress = ServerAddress
		}
		if clCfg.RequestTimeoutSec == 0 {
			clCfg.RequestTimeoutSec = RequestTimeoutSec
		}
		if clCfg.AdminAuthToken == "" {
			clCfg.AdminAuthToken = AdminAuthToken
		}
		if clCfg.SendMessageURI == "" {
			clCfg.SendMessageURI = clCfg.ServerAddress + SendMessageURL
		}
		if clCfg.StatusGetURI == "" {
			clCfg.StatusGetURI = clCfg.ServerAddress + GetStatusURL
		}
	}

	return clCfg, nil
}

// LoadClientConfigFile this method read a server configurations from a file in the json format.
func LoadClientConfigFile(configFilePath string) (*ClientConfig, error) {
	var configFromFile ClientConfig

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
