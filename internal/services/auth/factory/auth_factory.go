package factory

import (
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/antonioo83/notifier-server/internal/services/auth"
)

func NewUserAuthHandler(userRepository interfaces.UserRepository, config config.Config) *auth.UserAuthService {
	return auth.NewUserAuth(userRepository, config)
}
