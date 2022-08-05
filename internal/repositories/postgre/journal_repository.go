package postgre

import (
	"context"
	"github.com/antonioo83/notifier-server/internal/models"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4/pgxpool"
)

type journalRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewJournalRepository(context context.Context, pool *pgxpool.Pool) interfaces.JournalRepository {
	return &journalRepository{context, pool}
}

// Save creates a user in the database.
func (j journalRepository) Save(journal models.Journal) error {
	var lastInsertId int
	err := j.connection.QueryRow(
		j.context,
		`INSERT INTO ns_journal(message_id, user_id, resource_id, response_status, response_content, description) 
			     VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		&journal.MessageId, &journal.UserId, &journal.ResourceId, &journal.ResponseStatus, &journal.ResponseContent,
		&journal.Description,
	).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	return err
}
