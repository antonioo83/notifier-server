package auth

import (
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/models"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"net/http"
	"strings"
)

const (
	Admin = iota + 1
	User  = iota + 1
)

type UserAuth struct {
	Role int
	User models.User
}

type TokenProvider interface {
	FindByToken(code string) (*models.User, error)
}

type UserAuthService struct {
	tokenProvider TokenProvider
	config        config.Config
}

func NewUserAuth(userRepository interfaces.UserRepository, config config.Config) *UserAuthService {
	return &UserAuthService{userRepository, config}
}

func (u UserAuthService) GetAuthUser(token string) (*UserAuth, error) {
	if len(token) < 32 {
		return nil, fmt.Errorf("user doesn't have correct a token")
	}

	user, err := u.tokenProvider.FindByToken(token)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.AuthToken == u.config.Auth.AdminAuthToken && token == user.AuthToken {
		return &UserAuth{Role: Admin, User: *user}, nil
	}

	return &UserAuth{Role: User, User: *user}, nil
}

func (u UserAuthService) GetToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer", "", 1)

	return strings.ReplaceAll(token, " ", "")
}
