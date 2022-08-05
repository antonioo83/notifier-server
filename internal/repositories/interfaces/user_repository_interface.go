package interfaces

import "github.com/antonioo83/notifier-server/internal/models"

type UserRepository interface {
	// Save creates a user in the database.
	Save(user models.User) error
	// Update updates a user in the database.
	Update(model models.User) error
	// Delete deletes a user in the database.
	Delete(code string) error
	// FindByCode finds a user in the database by code.
	FindByCode(code string) (*models.User, error)
	// FindByToken finds a user in the database by token.
	FindByToken(code string) (*models.User, error)
	// FindAll finds users in the database by limit and offset.
	FindAll(limit int, offset int) (*map[int]models.User, error)
	// IsInDatabase user is exists in the database.
	IsInDatabase(code string) (bool, error)
}
