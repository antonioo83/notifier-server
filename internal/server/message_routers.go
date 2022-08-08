package server

import (
	"github.com/antonioo83/notifier-server/internal/handlers"
	"github.com/antonioo83/notifier-server/internal/services"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// getCreateMessagesRoute gets create messages route.
func getCreateMessagesRoute(r *chi.Mux, params services.MessageRouteParameters) *chi.Mux {
	r.Post("/api/v1/messages", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateMessageHandler(r, w, params)
	})

	return r
}

// getCreateMessagesRoute gets delete message route.
func getDeleteMessageRoute(r *chi.Mux, params services.MessageRouteParameters) *chi.Mux {
	r.Delete("/api/v1/messages", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeletedMessageHandler(r, w, params)
	})

	return r
}

// getMessageRoute gets get message route.
func getMessageRoute(r *chi.Mux, params services.MessageRouteParameters) *chi.Mux {
	r.Get("/api/v1/message", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMessageHandler(r, w, params)
	})

	return r
}
