package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antonioo83/notifier-server/internal/services"
	"github.com/antonioo83/notifier-server/internal/utils"
	"net/http"
	"strconv"
)

func GetCreatedUserResponse(r *http.Request, w http.ResponseWriter, param services.UserRouteParameters) {
	httpRequest, err := getRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authToken, err := services.CreateUser(*httpRequest, param)
	if err != nil {
		if errors.Is(err, services.RequestError) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsonResponse, err := getJSONResponse("token", authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

func getJSONResponse(key string, value string) ([]byte, error) {
	resp := make(map[string]string)
	resp[key] = value
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return jsonResp, fmt.Errorf("error happened in JSON marshal: %w", err)
	}

	return jsonResp, nil
}

func getRequest(r *http.Request) (*services.UserCreateRequest, error) {
	var request services.UserCreateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &request, nil
}

func GetUpdatedUserResponse(r *http.Request, w http.ResponseWriter, param services.UserRouteParameters) {
	httpRequest, err := getRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authToken, err := services.UpdateUser(*httpRequest, param)
	if err != nil {
		if errors.Is(err, services.RequestError) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	jsonResponse, err := getJSONResponse("token", authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

func GetDeletedUserResponse(r *http.Request, w http.ResponseWriter, param services.UserRouteParameters) {
	httpRequest, err := getDeleteRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = services.DeleteUser(*httpRequest, param)
	if err != nil {
		if errors.Is(err, services.RequestError) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func getDeleteRequest(r *http.Request) (*services.UserDeleteRequest, error) {
	var request services.UserDeleteRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &request, nil
}

func GetUserResponse(r *http.Request, w http.ResponseWriter, param services.UserRouteParameters) {
	httpRequest := getUserRequest(r)
	response, err := services.GetUser(*httpRequest, param)
	if err != nil {
		if errors.Is(err, services.RequestError) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if errors.Is(err, services.NotFoundError) {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	jsonResponse, err := getUserJsonResponse(*response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

func getUserRequest(r *http.Request) *services.UserGetRequest {
	var request services.UserGetRequest
	request.UserId = r.URL.Query().Get("userId")

	return &request
}

func getUserJsonResponse(resp services.UserResponse) ([]byte, error) {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return jsonResp, fmt.Errorf("error happened in JSON marshal: %w", err)
	}

	return jsonResp, nil
}

func GetUsersResponse(r *http.Request, w http.ResponseWriter, param services.UserRouteParameters) {
	httpRequest, err := getUsersRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responses, err := services.GetUsers(*httpRequest, param)
	if err != nil {
		if errors.Is(err, services.RequestError) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	jsonResponse, err := getUsersJsonResponse(responses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

func getUsersRequest(r *http.Request) (*services.UsersGetRequest, error) {
	var err error
	var request services.UsersGetRequest
	request.Limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		return nil, err
	}
	if request.Limit == 0 {
		request.Limit = 100
	}

	request.Offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func getUsersJsonResponse(resp []services.UserResponse) ([]byte, error) {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return jsonResp, fmt.Errorf("error happened in JSON marshal: %w", err)
	}

	return jsonResp, nil
}
