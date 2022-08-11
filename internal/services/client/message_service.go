package client

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/utils"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type messageService struct {
	cfg config.ClientConfig
}

func NewMessageService(cfg config.ClientConfig) MessageService {
	return &messageService{cfg}
}

type MessageRequest struct {
	MessageId string `validate:"required,max=64" json:"messageId,omitempty`
	URL       string `validate:"required,max=1000" json:"url,omitempty"`
	Command   string `validate:"required,max=10" json:"command,omitempty"`
	Content   string `validate:"required" json:"content,omitempty"`
}

func (m messageService) SendMessages(filepath string) (status int, err error) {
	messages, err := getMessagesFromFile(filepath)
	if err != nil {
		return 0, fmt.Errorf("i can't load file with messages: %w", err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * time.Duration(m.cfg.RequestTimeoutSec),
	}

	jsonReq, err := json.Marshal(messages)
	if err != nil {
		return 0, fmt.Errorf("i can't decode json request: %w", err)
	}

	req, err := http.NewRequest("POST", m.cfg.SendMessageURI, strings.NewReader(string(jsonReq)))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", m.cfg.AdminAuthToken)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	defer utils.ResourceClose(resp.Body)

	return resp.StatusCode, nil
}

// LoadFile this method read a server configurations from a file in the json format.
func getMessagesFromFile(filePath string) (*[]MessageRequest, error) {
	var messages []MessageRequest

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, fmt.Errorf("unable to open a file of messages: %w", err)
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

	err = json.Unmarshal(jsonConfig, &messages)
	if err != nil {
		return nil, fmt.Errorf("i can't parse a configuration json file: %w", err)
	}

	return &messages, nil
}

type MessageStatusResponse struct {
	MessageId string `json:"messageId`
	IsSent    bool   `json:"isSent"`
}

func (m messageService) GetStatus(messageId string) (MessageStatusResponse, error) {
	var result MessageStatusResponse
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * time.Duration(m.cfg.RequestTimeoutSec),
	}

	req, err := http.NewRequest("GET", m.cfg.StatusGetURI, nil)
	if err != nil {
		return result, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", m.cfg.AdminAuthToken)
	values := req.URL.Query()
	values.Add("messageId", messageId)
	req.URL.RawQuery = values.Encode()

	response, err := client.Do(req)
	if err != nil {
		return result, err
	}
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	defer utils.ResourceClose(response.Body)

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
