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
	"strings"
	"testing"
)

type SettingRequestTest struct {
	url             string
	queryParams     map[string]string
	method          string
	content         string
	responseStatus  int
	responseContent string
	description     string
}

// TestGetRouters test crud actions for users.
func TestGetSettingRouters(t *testing.T) {
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

	repository := factory.NewSettingRepository(context, pool)
	userRepository := factory.NewUserRepository(context, pool)
	userAuthHandler := factory2.NewUserAuthHandler(userRepository, cfg)
	routeParameters :=
		RouteParameters{
			Config:            cfg,
			SettingRepository: repository,
		}

	r := GetRouters(userAuthHandler, routeParameters)
	ts := httptest.NewServer(r)
	defer ts.Close()

	var tests []SettingRequestTest
	createRequest := services.SettingCreateRequest{}
	err = faker.FakeData(&createRequest)
	assert.NoError(t, err)
	createContent, err := json.Marshal(createRequest)
	assert.NoError(t, err)
	tests = append(tests, SettingRequestTest{
		url:            ts.URL + "/api/v1/defaultSettings",
		method:         "POST",
		content:        string(createContent),
		responseStatus: 201,
		description:    "setting create",
	})

	afterCreateResponse := services.SettingResponse{}
	err = copier.Copy(&afterCreateResponse, &createRequest)
	assert.NoError(t, err)
	afterCreateResponse.SettingId = createRequest.SettingId
	afterCreateResponseJson, err := json.Marshal(afterCreateResponse)
	assert.NoError(t, err)
	q := make(map[string]string)
	q["settingId"] = createRequest.SettingId
	tests = append(tests, SettingRequestTest{
		url:             ts.URL + "/api/v1/defaultSetting",
		queryParams:     q,
		method:          "GET",
		responseStatus:  200,
		responseContent: string(afterCreateResponseJson),
		description:     "get created setting",
	})

	updateRequest := services.SettingCreateRequest{}
	err = faker.FakeData(&updateRequest)
	assert.NoError(t, err)
	updateRequest.SettingId = createRequest.SettingId
	updateContent, err := json.Marshal(updateRequest)
	assert.NoError(t, err)
	tests = append(tests, SettingRequestTest{
		url:            ts.URL + "/api/v1/defaultSettings",
		method:         "PUT",
		content:        string(updateContent),
		responseStatus: 204,
		description:    "setting update",
	})

	afterUpdateResponse := services.SettingResponse{}
	err = copier.Copy(&afterUpdateResponse, &updateRequest)
	assert.NoError(t, err)
	afterUpdateResponse.SettingId = createRequest.SettingId
	afterUpdateResponseJson, err := json.Marshal(afterUpdateResponse)
	assert.NoError(t, err)
	tests = append(tests, SettingRequestTest{
		url:             ts.URL + "/api/v1/defaultSetting",
		queryParams:     q,
		method:          "GET",
		responseStatus:  200,
		responseContent: string(afterUpdateResponseJson),
		description:     "get updated setting",
	})

	deleteRequest := services.SettingDeleteRequest{SettingId: createRequest.SettingId}
	deleteContent, err := json.Marshal(deleteRequest)
	assert.NoError(t, err)
	tests = append(tests, SettingRequestTest{
		url:            ts.URL + "/api/v1/defaultSettings",
		method:         "DELETE",
		content:        string(deleteContent),
		responseStatus: 202,
		description:    "setting delete",
	})

	assert.NoError(t, err)
	tests = append(tests, SettingRequestTest{
		url:            ts.URL + "/api/v1/defaultSetting",
		queryParams:    q,
		method:         "GET",
		responseStatus: 204,
		description:    "get deleted setting",
	})

	for _, t1 := range tests {
		fmt.Println("start test:" + t1.description)
		createRequest, err := getRequest(t1.url, t1.queryParams, t1.method, strings.NewReader(t1.content), cfg.Auth.AdminAuthToken)
		assert.NoError(t, err)
		resp, respBody := sendRequest(t, createRequest)
		require.NoError(t, err)
		assert.Equal(t, t1.responseStatus, resp.StatusCode)
		if t1.responseContent != "" {
			assert.Equal(t, t1.responseContent, respBody)
		}
		assert.NoError(t, err)
	}
}
