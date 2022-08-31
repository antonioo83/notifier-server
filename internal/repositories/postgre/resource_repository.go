package postgre

import (
	"context"
	"errors"
	"github.com/antonioo83/notifier-server/internal/models"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type resourceRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewResourceRepository(context context.Context, pool *pgxpool.Pool) interfaces.ResourceRepository {
	return &resourceRepository{context, pool}
}

// Save creates a resource in the database.
func (r resourceRepository) Save(resource models.Resource) (int, error) {
	var lastInsertId int
	err := r.connection.QueryRow(
		r.context,
		"INSERT INTO ns_resources(code, user_id, url, description)VALUES ($1, $2, $3, $4) RETURNING id",
		&resource.Code, &resource.UserId, &resource.URL, &resource.Description,
	).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}

	return lastInsertId, err
}

// Delete deletes a resource from the database.
func (r resourceRepository) Delete(code string) error {
	_, err := r.connection.Exec(r.context, "UPDATE ns_resources SET deleted_at=NOW() WHERE code=$1 AND deleted_at IS NULL", code)

	return err
}

// FindByCode find a resource by code.
func (r resourceRepository) FindByCode(code string) (*models.Resource, error) {
	var model models.Resource
	err := r.connection.QueryRow(
		r.context,
		"SELECT id, code, user_id, url, description, created_at FROM ns_resources WHERE code=$1 AND deleted_at IS NULL",
		code,
	).Scan(&model.ID, &model.Code, &model.UserId, &model.URL, &model.Description, &model.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

// IsInDatabase a resource exists in the database.
func (r resourceRepository) IsInDatabase(code string) (bool, error) {
	model, err := r.FindByCode(code)

	return !(model == nil), err
}
