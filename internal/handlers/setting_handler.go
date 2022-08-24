package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antonioo83/notifier-server/internal/services"
	"github.com/antonioo83/notifier-server/internal/services/auth"
	"github.com/antonioo83/notifier-server/internal/utils"
	"net/http"
	"strconv"
)

// CreateSettingHandler user create handler.
func CreateSettingHandler(r *http.Request, w http.ResponseWriter, param services.SettingRouteParameters) {
	httpRequest, err := getCreateSettingRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value("userAuth").(*auth.UserAuth)
	err = services.CreateSetting(*userAuth, *httpRequest, param)
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
}

// getCreateSettingRequest returns request by json request.
func getCreateSettingRequest(r *http.Request) (*services.SettingCreateRequest, error) {
	var request services.SettingCreateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &request, nil
}

// UpdateSettingHandler user update handler.
func UpdateSettingHandler(r *http.Request, w http.ResponseWriter, param services.SettingRouteParameters) {
	httpRequest, err := getCreateSettingRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value("userAuth").(*auth.UserAuth)
	err = services.UpdateSetting(*userAuth, *httpRequest, param)
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
}

// DeleteSettingHandler user delete handler.
func DeleteSettingHandler(r *http.Request, w http.ResponseWriter, param services.SettingRouteParameters) {
	httpRequest, err := getDeleteSettingRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value("userAuth").(*auth.UserAuth)
	err = services.DeleteSetting(*userAuth, *httpRequest, param)
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

// getDeleteRequest returns request for the delete a user.
func getDeleteSettingRequest(r *http.Request) (*services.SettingDeleteRequest, error) {
	var request services.SettingDeleteRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &request, nil
}

// GetSettingHandler user get handler.
func GetSettingHandler(r *http.Request, w http.ResponseWriter, param services.SettingRouteParameters) {
	httpRequest := getSettingRequest(r)
	userAuth := r.Context().Value("userAuth").(*auth.UserAuth)
	response, err := services.GetSetting(*userAuth, *httpRequest, param)
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
	jsonResponse, err := getSettingJsonResponse(*response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

// getSettingRequest returns request for get a user.
func getSettingRequest(r *http.Request) *services.SettingGetRequest {
	var request services.SettingGetRequest
	request.SettingId = r.URL.Query().Get("settingId")

	return &request
}

// getUserJsonResponse returns a response in json format.
func getSettingJsonResponse(resp services.SettingResponse) ([]byte, error) {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return jsonResp, fmt.Errorf("error happened in JSON marshal: %w", err)
	}

	return jsonResp, nil
}

// GetSettingsHandler users get handler.
func GetSettingsHandler(r *http.Request, w http.ResponseWriter, param services.SettingRouteParameters) {
	httpRequest, err := getSettingsRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value("userAuth").(*auth.UserAuth)
	responses, err := services.GetSettings(*userAuth, *httpRequest, param)
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
	jsonResponse, err := getSettingsJsonResponse(responses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

// getSettingsRequest returns request for get users.
func getSettingsRequest(r *http.Request) (*services.SettingsGetRequest, error) {
	var err error
	var request services.SettingsGetRequest
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

// getSettingsJsonResponse returns a users response in json format.
func getSettingsJsonResponse(resp []services.SettingResponse) ([]byte, error) {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return jsonResp, fmt.Errorf("error happened in JSON marshal: %w", err)
	}

	return jsonResp, nil
}
