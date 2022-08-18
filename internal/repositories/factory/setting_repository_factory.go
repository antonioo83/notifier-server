package factory

import (
	"context"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/antonioo83/notifier-server/internal/repositories/postgre"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewSettingRepository(context context.Context, pool *pgxpool.Pool) interfaces.SettingRepository {
	return postgre.NewSettingRepository(context, pool)
}
