package interfaces

import "github.com/antonioo83/notifier-server/internal/models"

type MessageRepository interface {
	Save(user models.Message) error
	Update(model models.Message) error
	Delete(code string) error
	FindByCode(code string) (*models.Message, error)
	FindAll(attemptCountMax int, limit int, offset int) (*map[int]models.Message, error)
	MarkSent(code string) error
	MarkUnSent(code string, attemptCount int) error
	IsInDatabase(code string) (bool, error)
}
