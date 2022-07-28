package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/repositories/factory"
	"github.com/antonioo83/notifier-server/internal/services"
	factory2 "github.com/antonioo83/notifier-server/internal/services/auth/factory"
	"github.com/antonioo83/notifier-server/internal/utils"
	"github.com/bxcodec/faker/v3"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type RequestTest struct {
	url             string
	queryParams     map[string]string
	method          string
	content         string
	responseStatus  int
	responseContent string
	description     string
}

func TestGetRouters(t *testing.T) {
	var pool *pgxpool.Pool
	context := context.Background()
	config, err := config.GetConfigSettings()
	if err != nil {
		log.Fatalf("Can't read config: %s", err.Error())
	}

	pool, _ = pgxpool.Connect(context, config.DatabaseDsn)
	defer pool.Close()

	userRepository := factory.NewUserRepository(context, pool)
	userAuthHandler := factory2.NewUserAuthHandler(userRepository, config)
	routeParameters :=
		RouteParameters{
			Config:         config,
			UserRepository: userRepository,
		}

	r := GetRouters(userAuthHandler, routeParameters)
	ts := httptest.NewServer(r)
	defer ts.Close()

	var tests []RequestTest
	createRequest := services.UserCreateRequest{}
	err = faker.FakeData(&createRequest)
	assert.NoError(t, err)
	createContent, err := json.Marshal(createRequest)
	assert.NoError(t, err)
	tests = append(tests, RequestTest{
		url:            ts.URL + "/api/v1/users",
		method:         "POST",
		content:        string(createContent),
		responseStatus: 201,
		description:    "user create",
	})

	afterCreateResponse := services.UserResponse{}
	err = copier.Copy(&afterCreateResponse, &createRequest)
	assert.NoError(t, err)
	afterCreateResponse.UserId = createRequest.UserId
	afterCreateResponseJson, err := json.Marshal(afterCreateResponse)
	assert.NoError(t, err)
	q := make(map[string]string)
	q["userId"] = createRequest.UserId
	tests = append(tests, RequestTest{
		url:             ts.URL + "/api/v1/user",
		queryParams:     q,
		method:          "GET",
		responseStatus:  200,
		responseContent: string(afterCreateResponseJson),
		description:     "get created user",
	})

	updateRequest := services.UserCreateRequest{}
	err = faker.FakeData(&updateRequest)
	assert.NoError(t, err)
	updateRequest.UserId = createRequest.UserId
	updateContent, err := json.Marshal(updateRequest)
	assert.NoError(t, err)
	tests = append(tests, RequestTest{
		url:            ts.URL + "/api/v1/users",
		method:         "PUT",
		content:        string(updateContent),
		responseStatus: 204,
		description:    "user update",
	})

	afterUpdateResponse := services.UserResponse{}
	err = copier.Copy(&afterUpdateResponse, &updateRequest)
	assert.NoError(t, err)
	afterUpdateResponse.UserId = createRequest.UserId
	afterUpdateResponseJson, err := json.Marshal(afterUpdateResponse)
	assert.NoError(t, err)
	tests = append(tests, RequestTest{
		url:             ts.URL + "/api/v1/user",
		queryParams:     q,
		method:          "GET",
		responseStatus:  200,
		responseContent: string(afterUpdateResponseJson),
		description:     "get updated user",
	})

	deleteRequest := services.UserDeleteRequest{UserId: createRequest.UserId}
	deleteContent, err := json.Marshal(deleteRequest)
	assert.NoError(t, err)
	tests = append(tests, RequestTest{
		url:            ts.URL + "/api/v1/users",
		method:         "DELETE",
		content:        string(deleteContent),
		responseStatus: 202,
		description:    "user delete",
	})

	assert.NoError(t, err)
	tests = append(tests, RequestTest{
		url:            ts.URL + "/api/v1/user",
		queryParams:    q,
		method:         "GET",
		responseStatus: 204,
		description:    "get deleted user",
	})

	for _, t1 := range tests {
		fmt.Println("start test:" + t1.description)
		createRequest, err := getRequest(t1.url, t1.queryParams, t1.method, strings.NewReader(t1.content), config.Auth.AdminAuthToken)
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

func getRequest(url string, queryParameters map[string]string, method string, body io.Reader, token string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	if queryParameters != nil {
		values := req.URL.Query()
		for param, value := range queryParameters {
			values.Add(param, value)
		}
		req.URL.RawQuery = values.Encode()
	}

	return req, err
}

func sendRequest(t *testing.T, req *http.Request) (*http.Response, string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer utils.ResourceClose(resp.Body)

	return resp, string(respBody)
}
