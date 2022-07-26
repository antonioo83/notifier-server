package interfaces

import "github.com/antonioo83/notifier-server/internal/models"

type UserRepository interface {
	Save(user models.User) error
	Update(model models.User) error
	Delete(code string) error
	FindByCode(code string) (*models.User, error)
	FindByToken(code string) (*models.User, error)
	FindAll(limit int, offset int) (*map[int]models.User, error)
	IsInDatabase(code string) (bool, error)
}
