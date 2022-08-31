package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/repositories/factory"
	"github.com/antonioo83/notifier-server/internal/services"
	factory2 "github.com/antonioo83/notifier-server/internal/services/auth/factory"
	"github.com/bxcodec/faker/v3"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

type MessageRequestTest struct {
	url             string
	queryParams     map[string]string
	method          string
	content         string
	responseStatus  int
	responseContent string
	description     string
}

// TestGetMessageRouters test crud actions for messages.
func TestGetMessageRouters(t *testing.T) {
	var pool *pgxpool.Pool
	context := context.Background()
	configPath, err := GetConfigPath()
	if err != nil {
		log.Fatalf("i can't get path to the configuration file:" + err.Error())
	}
	configFromFile, err := config.LoadConfigFile(configPath)
	if err != nil {
		log.Fatalf("i can't load configuration file:" + err.Error())
	}
	cfg, err := config.GetConfigSettings(configFromFile)
	if err != nil {
		log.Fatalf("Can't read config: %s", err.Error())
	}

	pool, _ = pgxpool.Connect(context, cfg.DatabaseDsn)
	defer pool.Close()

	userRepository := factory.NewUserRepository(context, pool)
	settingRepository := factory.NewSettingRepository(context, pool)
	resourceRepository := factory.NewResourceRepository(context, pool)
	messageRepository := factory.NewMessageRepository(context, pool)
	userAuthHandler := factory2.NewUserAuthHandler(userRepository, cfg)
	journalRepository := factory.NewJournalRepository(context, pool)
	routeParameters :=
		RouteParameters{
			Config:             cfg,
			UserRepository:     userRepository,
			SettingRepository:  settingRepository,
			ResourceRepository: resourceRepository,
			MessageRepository:  messageRepository,
			JournalRepository:  journalRepository,
		}

	r := GetRouters(userAuthHandler, routeParameters)
	ts := httptest.NewServer(r)
	defer ts.Close()

	var tests []MessageRequestTest
	var createRequests []services.MessageCreateRequest
	createRequest := services.MessageCreateRequest{}
	err = faker.FakeData(&createRequest)
	assert.NoError(t, err)
	createRequests = append(createRequests, createRequest)
	createContent, err := json.Marshal(createRequests)
	assert.NoError(t, err)
	tests = append(tests, MessageRequestTest{
		url:            ts.URL + "/api/v1/messages",
		method:         "POST",
		content:        string(createContent),
		responseStatus: 201,
		description:    "message create",
	})

	afterCreateResponse := services.MessageResponse{}
	err = copier.Copy(&afterCreateResponse, &createRequest)
	assert.NoError(t, err)
	afterCreateResponse.MessageId = createRequest.MessageId
	afterCreateResponseJson, err := json.Marshal(afterCreateResponse)
	assert.NoError(t, err)
	q := make(map[string]string)
	q["messageId"] = createRequest.MessageId
	tests = append(tests, MessageRequestTest{
		url:             ts.URL + "/api/v1/message",
		queryParams:     q,
		method:          "GET",
		responseStatus:  200,
		responseContent: string(afterCreateResponseJson),
		description:     "get created message",
	})

	deleteRequest := services.MessageDeleteRequest{MessageId: createRequest.MessageId}
	deleteContent, err := json.Marshal(deleteRequest)
	assert.NoError(t, err)
	tests = append(tests, MessageRequestTest{
		url:            ts.URL + "/api/v1/messages",
		method:         "DELETE",
		content:        string(deleteContent),
		responseStatus: 202,
		description:    "message delete",
	})

	assert.NoError(t, err)
	tests = append(tests, MessageRequestTest{
		url:            ts.URL + "/api/v1/message",
		queryParams:    q,
		method:         "GET",
		responseStatus: 204,
		description:    "get deleted message",
	})

	for _, t1 := range tests {
		fmt.Println("start test:" + t1.description)
		req, err := getRequest(t1.url, t1.queryParams, t1.method, strings.NewReader(t1.content), cfg.Auth.AdminAuthToken)
		assert.NoError(t, err)
		resp, respBody := sendRequest(t, req)
		require.NoError(t, err)

		if t1.method == "GET" {
			m := regexp.MustCompile(`("createdAt":"[0-9]{4}-[0-9]{2}-[0-9]{2}\s[0-9]{2}:[0-9]{2}:[0-9]{2}")`)
			respBody = m.ReplaceAllString(respBody, `"createdAt":""`)
		}

		assert.Equal(t, t1.responseStatus, resp.StatusCode)
		if t1.responseContent != "" {
			assert.Equal(t, t1.responseContent, respBody)
		}
		assert.NoError(t, err)
	}
}
