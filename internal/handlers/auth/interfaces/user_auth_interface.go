package interfaces

import (
	"github.com/antonioo83/notifier-server/internal/handlers/auth"
	"net/http"
)

type UserAuthHandler interface {
	GetAuthUser(token string) (auth.UserAuth, error)
	GetToken(r *http.Request) string
}
