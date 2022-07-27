package main

import (
	"context"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/repositories/factory"
	"github.com/antonioo83/notifier-server/internal/server"
	factory2 "github.com/antonioo83/notifier-server/internal/services/auth/factory"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

func main() {
	config, err := config.GetConfigSettings()
	if err != nil {
		log.Fatalf("Can't resd config: %s", err.Error())
	}

	var pool *pgxpool.Pool
	context := context.Background()
	pool, err = pgxpool.Connect(context, config.DatabaseDsn)
	if err != nil {
		log.Fatalf("Can't connect to the database server: %s", err.Error())
	}
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
