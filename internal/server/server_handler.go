package server

import (
	"compress/flate"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/handlers"
	"github.com/antonioo83/notifier-server/internal/handlers/auth"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

type RouteParameters struct {
	Config         config.Config
	UserRepository interfaces.UserRepository
}

func GetRouters(uh *auth.UserAuthHandler, p RouteParameters) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(p.Config.RequestTimeoutSec) * time.Second))
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(compressor.Handler)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := uh.GetToken(r)
			userAuth, err := uh.GetAuthUser(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if strings.HasPrefix(r.RequestURI, "/api/v1/users") == true && userAuth.Role != auth.Admin {
				http.Error(w, "access for this route is denied", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), "userAuth", userAuth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	var params = handlers.UserRouteParameters{
		Config:         p.Config,
		UserRepository: p.UserRepository,
	}
	r = getCreateUserRoute(r, params)
	r = getUpdateUserRoute(r, params)
	r = getDeleteUserRoute(r, params)
	r = getUserRoute(r, params)
	r = getUsersRoute(r, params)

	return r
}
