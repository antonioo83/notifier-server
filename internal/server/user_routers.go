package server

import (
	"github.com/antonioo83/notifier-server/internal/handlers"
	"github.com/antonioo83/notifier-server/internal/services"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func getCreateUserRoute(r *chi.Mux, params services.UserRouteParameters) *chi.Mux {
	r.Post("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatedUserHandler(r, w, params)
	})

	return r
}

func getUpdateUserRoute(r *chi.Mux, params services.UserRouteParameters) *chi.Mux {
	r.Put("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdatedUserHandler(r, w, params)
	})

	return r
}

func getDeleteUserRoute(r *chi.Mux, params services.UserRouteParameters) *chi.Mux {
	r.Delete("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeletedUserHandler(r, w, params)
	})

	return r
}

func getUserRoute(r *chi.Mux, params services.UserRouteParameters) *chi.Mux {
	r.Get("/api/v1/user", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserHandler(r, w, params)
	})

	return r
}

func getUsersRoute(r *chi.Mux, params services.UserRouteParameters) *chi.Mux {
	r.Get("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUsersHandler(r, w, params)
	})

	return r
}
