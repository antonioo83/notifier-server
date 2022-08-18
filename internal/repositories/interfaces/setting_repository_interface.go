package interfaces

import "github.com/antonioo83/notifier-server/internal/models"

type SettingRepository interface {
	// Save creates a user in the database.
	Save(model models.Setting) error
	// Update updates a user in the database.
	Update(model models.Setting) error
	// Delete deletes a user in the database.
	Delete(resourceId string) error
	// FindByCode finds a user in the database by code.
	FindByCode(resourceId string) (*models.Setting, error)
	// FindAll finds users in the database by limit and offset.
	FindAll(limit int, offset int) (*map[int]models.Setting, error)
	// IsInDatabase user is exists in the database.
	IsInDatabase(resourceId string) (bool, error)
}
