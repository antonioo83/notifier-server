package server

import (
	"github.com/antonioo83/notifier-server/internal/handlers"
	"github.com/antonioo83/notifier-server/internal/services"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// getCreateSettingRoute get create user route.
func getCreateSettingRoute(r *chi.Mux, params services.SettingRouteParameters) *chi.Mux {
	r.Post("/api/v1/defaultSettings", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateSettingHandler(r, w, params)
	})

	return r
}

// getCreateSettingRoute get update user route.
func getUpdateSettingRoute(r *chi.Mux, params services.SettingRouteParameters) *chi.Mux {
	r.Put("/api/v1/defaultSettings", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateSettingHandler(r, w, params)
	})

	return r
}

// getCreateSettingRoute get delete user route.
func getDeleteSettingRoute(r *chi.Mux, params services.SettingRouteParameters) *chi.Mux {
	r.Delete("/api/v1/defaultSettings", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteSettingHandler(r, w, params)
	})

	return r
}

// getCreateSettingRoute get an user route.
func getSettingRoute(r *chi.Mux, params services.SettingRouteParameters) *chi.Mux {
	r.Get("/api/v1/defaultSetting", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetSettingHandler(r, w, params)
	})

	return r
}

// getCreateSettingRoute get users route.
func getSettingsRoute(r *chi.Mux, params services.SettingRouteParameters) *chi.Mux {
	r.Get("/api/v1/defaultSettings", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetSettingsHandler(r, w, params)
	})

	return r
}
