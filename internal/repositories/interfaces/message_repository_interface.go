package interfaces

import "github.com/antonioo83/notifier-server/internal/models"

type MessageRepository interface {
	// Save create a message in the database.
	Save(user models.Message) error
	// Update modify a message in the database.
	Update(model models.Message) error
	// Delete delete a message in the database.
	Delete(code string) error
	// FindByCode find a message by code.
	FindByCode(code string) (*models.Message, error)
	// FindAll find messages by limit and offset.
	FindAll(attemptCountMax int, limit int, offset int) (*map[int]models.Message, error)
	// MarkSent mark a message as sent.
	MarkSent(code string) error
	// MarkUnSent mark a message as not sent.
	MarkUnSent(code string, attemptCount int) error
	// IsInDatabase a message exists in the database.
	IsInDatabase(code string) (bool, error)
}
