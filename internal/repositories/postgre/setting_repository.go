package postgre

import (
	"context"
	"errors"
	"github.com/antonioo83/notifier-server/internal/models"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type settingRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewSettingRepository(context context.Context, pool *pgxpool.Pool) interfaces.SettingRepository {
	return &settingRepository{context, pool}
}

// Save creates a user in the database.
func (s settingRepository) Save(model models.Setting) error {
	var lastInsertId int
	err := s.connection.QueryRow(
		s.context,
		`INSERT INTO ns_settings (code, user_id, resource_id, title, callback_url, count, intervals, timeout, description)
               VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		&model.Code, &model.UserId, &model.ResourceId, &model.Title, &model.CallbackURL, &model.Count, &model.Intervals,
		&model.Timeout, &model.Description,
	).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	return err
}

// Update updates a user in the database.
func (s settingRepository) Update(model models.Setting) error {
	_, err := s.connection.Exec(
		s.context,
		`UPDATE 
               ns_settings 
             SET 
               title=$1, callback_url=$2, count=$3, intervals=$4, description=$5, updated_at=NOW() 
             WHERE 
               user_id=$6 AND resource_id=&7 AND code=$8 AND deleted_at IS NULL`,
		&model.Title, &model.CallbackURL, &model.Count, &model.Intervals, &model.Description,
		&model.UserId, &model.ResourceId, &model.Code,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a user in the database.
func (s settingRepository) Delete(resourceId string) error {
	_, err := s.connection.Exec(
		s.context,
		"UPDATE ns_settings SET deleted_at=NOW() WHERE resource_id=&1 AND deleted_at IS NULL",
		&resourceId,
	)

	return err
}

// FindByCode finds a user in the database by code.
func (s settingRepository) FindByCode(resourceId string) (*models.Setting, error) {
	var model models.Setting
	err := s.connection.QueryRow(
		s.context,
		`SELECT 
               id, code, user_id, resource_id, title, callback_url, count, intervals, timeout, description, created_at 
             FROM 
               ns_settings 
             WHERE 
               resource_id=$1 AND deleted_at IS NULL`,
		resourceId,
	).Scan(&model.ID, &model.Code, &model.UserId, &model.ResourceId, &model.Title, &model.CallbackURL, &model.Count,
		&model.Intervals, &model.Timeout, &model.Description, &model.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

// FindAll finds users in the database by limit and offset.
func (s settingRepository) FindAll(limit int, offset int) (*map[int]models.Setting, error) {
	rows, err := s.connection.Query(
		s.context,
		`SELECT 
				id, code, user_id, resource_id, title, callback_url, count, intervals, timeout, description, created_at
			FROM ns_settings s
			WHERE 
  				s.deleted_at IS NULL 
			ORDER BY s.id ASC
			LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}

	models, err := getSettingModels(rows)
	if err != nil {
		return nil, err
	}

	return &models, nil
}

// getModels returns array of the user models.
func getSettingModels(rows pgx.Rows) (map[int]models.Setting, error) {
	var settings = make(map[int]models.Setting)
	var model models.Setting
	for rows.Next() {
		err := rows.Scan(
			&model.ID, &model.Code, &model.UserId, &model.ResourceId, &model.Title, &model.CallbackURL, &model.Count,
			&model.Intervals, &model.Timeout, &model.Description, &model.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		settings[model.ID] = model
	}

	return settings, nil
}

// IsInDatabase user is exists in the database.
func (s settingRepository) IsInDatabase(resourceId string) (bool, error) {
	model, err := s.FindByCode(resourceId)

	return !(model == nil), err
}
