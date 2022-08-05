package interfaces

import "github.com/antonioo83/notifier-server/internal/models"

type ResourceRepository interface {
	// Save creates a resource in the database.
	Save(model models.Resource) (int, error)
	// Delete deletes a resource from the database.
	Delete(code int) error
	// FindByCode find a resource by code.
	FindByCode(code int) (*models.Resource, error)
	// IsInDatabase a resource exists in the database.
	IsInDatabase(code int) (bool, error)
}
