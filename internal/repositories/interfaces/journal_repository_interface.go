package interfaces

import "github.com/antonioo83/notifier-server/internal/models"

type JournalRepository interface {
	// Save create journal record.
	Save(journal models.Journal) error
}
