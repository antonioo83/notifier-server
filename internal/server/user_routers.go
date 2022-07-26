package server

import (
	"github.com/antonioo83/notifier-server/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func getCreateUserRoute(r *chi.Mux, params handlers.UserRouteParameters) *chi.Mux {
	r.Post("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreatedUserResponse(r, w, params)
	})

	return r
}

func getUpdateUserRoute(r *chi.Mux, params handlers.UserRouteParameters) *chi.Mux {
	r.Put("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUpdatedUserResponse(r, w, params)
	})

	return r
}

func getDeleteUserRoute(r *chi.Mux, params handlers.UserRouteParameters) *chi.Mux {
	r.Delete("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetDeletedUserResponse(r, w, params)
	})

	return r
}

func getUserRoute(r *chi.Mux, params handlers.UserRouteParameters) *chi.Mux {
	r.Get("/api/v1/user", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserResponse(r, w, params)
	})

	return r
}

func getUsersRoute(r *chi.Mux, params handlers.UserRouteParameters) *chi.Mux {
	r.Get("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUsersResponse(r, w, params)
	})

	return r
}
