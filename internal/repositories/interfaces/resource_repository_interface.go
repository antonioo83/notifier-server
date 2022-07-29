package interfaces

import "github.com/antonioo83/notifier-server/internal/models"

type ResourceRepository interface {
	Save(model models.Resource) (int, error)
	Delete(code int) error
	FindByCode(code int) (*models.Resource, error)
	IsInDatabase(code int) (bool, error)
}
