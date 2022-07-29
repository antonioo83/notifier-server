package factory

import (
	"context"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/antonioo83/notifier-server/internal/repositories/postgre"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewMessageRepository(context context.Context, pool *pgxpool.Pool) interfaces.MessageRepository {
	return postgre.NewMessageRepository(context, pool)
}
