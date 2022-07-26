package main

import (
	"context"
	"github.com/antonioo83/notifier-server/config"
	factory2 "github.com/antonioo83/notifier-server/internal/handlers/auth/factory"
	"github.com/antonioo83/notifier-server/internal/repositories/factory"
	"github.com/antonioo83/notifier-server/internal/server"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

func main() {
	config := config.GetConfigSettings()
	var pool *pgxpool.Pool
	context := context.Background()
	pool, _ = pgxpool.Connect(context, config.DatabaseDsn)
	defer pool.Close()

	userRepository := factory.NewUserRepository(context, pool)
	userAuthHandler := factory2.NewUserAuthHandler(userRepository, config)
	routeParameters :=
		server.RouteParameters{
			Config:         config,
			UserRepository: userRepository,
		}

	handler := server.GetRouters(userAuthHandler, routeParameters)
	log.Fatal(http.ListenAndServe(config.ServerAddress, handler))
}
