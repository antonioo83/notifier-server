package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/notifier-server/internal/services"
	"github.com/antonioo83/notifier-server/internal/services/auth"
	"github.com/antonioo83/notifier-server/internal/utils"
	"net/http"
)

func GetCreateMessageHandler(r *http.Request, w http.ResponseWriter, param services.MessageRouteParameters) {
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

func getMessageRequests(r *http.Request) (*[]services.MessageCreateRequest, error) {
	var requests []services.MessageCreateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requests)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &requests, nil
}
