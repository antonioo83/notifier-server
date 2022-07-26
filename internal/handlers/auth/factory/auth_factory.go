package factory

import (
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/handlers/auth"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
)

func NewUserAuthHandler(userRepository interfaces.UserRepository, config config.Config) *auth.UserAuthHandler {
	return auth.NewUserAuth(userRepository, config)
}
