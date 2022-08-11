package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antonioo83/notifier-server/internal/services"
	"github.com/antonioo83/notifier-server/internal/services/auth"
	"github.com/antonioo83/notifier-server/internal/utils"
	"net/http"
)

// CreateMessageHandler create a message handler.
func CreateMessageHandler(r *http.Request, w http.ResponseWriter, param services.MessageRouteParameters) {
	httpRequests, err := getMessageRequests(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value("userAuth").(*auth.UserAuth)
	responses, err := services.CreateMessages(*userAuth, *httpRequests, param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

// getMessageRequests returns a request for create a message.
func getMessageRequests(r *http.Request) (*[]services.MessageCreateRequest, error) {
	var requests []services.MessageCreateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requests)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &requests, nil
}

// DeletedMessageHandler Delete message handler.
func DeletedMessageHandler(r *http.Request, w http.ResponseWriter, param services.MessageRouteParameters) {
	httpRequest, err := getMessageDeleteRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = services.DeleteMessage(*httpRequest, param)
	if err != nil {
		if errors.Is(err, services.RequestError) {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

// getMessageDeleteRequest returns request for delete a message.
func getMessageDeleteRequest(r *http.Request) (*services.MessageDeleteRequest, error) {
	var request services.MessageDeleteRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &request, nil
}

// GetMessageHandler get message handler.
func GetMessageHandler(r *http.Request, w http.ResponseWriter, param services.MessageRouteParameters) {
	httpRequest := getMessageRequest(r)
	response, err := services.GetMessage(*httpRequest, param)
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
	jsonResponse, err := getMessageJsonResponse(*response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

// getMessageRequest returns request for get a message.
func getMessageRequest(r *http.Request) *services.MessageGetRequest {
	var request services.MessageGetRequest
	request.MessageId = r.URL.Query().Get("messageId")

	return &request
}

// getMessageJsonResponse converts a response to the json response.
func getMessageJsonResponse(resp services.MessageResponse) ([]byte, error) {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return jsonResp, fmt.Errorf("error happened in JSON marshal: %w", err)
	}

	return jsonResp, nil
}
