package server

import (
	"github.com/antonioo83/notifier-server/internal/handlers"
	"github.com/antonioo83/notifier-server/internal/services"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func getCreateMessagesRoute(r *chi.Mux, params services.MessageRouteParameters) *chi.Mux {
	r.Post("/api/v1/messages", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreateMessageHandler(r, w, params)
	})

	return r
}

func getDeleteMessageRoute(r *chi.Mux, params services.UserRouteParameters) *chi.Mux {
	r.Delete("/api/v1/messages", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetDeletedUserResponse(r, w, params)
	})

	return r
}

func getMessageRoute(r *chi.Mux, params services.UserRouteParameters) *chi.Mux {
	r.Get("/api/v1/message", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserResponse(r, w, params)
	})

	return r
}
