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

// Save creates a setting in the database.
func (s settingRepository) Save(model models.Setting) error {
	var lastInsertId int
	err := s.connection.QueryRow(
		s.context,
		`INSERT INTO ns_settings (code, user_id, resource_id, title, callback_url, count, intervals, timeout, description)
               VALUES ($1, $2, NULLIF($3,0), $4, $5, $6, $7, $8, $9) RETURNING id`,
		&model.Code, &model.UserId, &model.ResourceId, &model.Title, &model.CallbackURL, &model.Count, &model.Intervals,
		&model.Timeout, &model.Description,
	).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	return err
}

// Update updates a setting in the database.
func (s settingRepository) Update(model models.Setting) error {
	_, err := s.connection.Exec(
		s.context,
		`UPDATE 
               ns_settings 
             SET 
               title=$1, callback_url=$2, count=$3, intervals=$4, timeout=$5, description=$6, updated_at=NOW() 
             WHERE 
               user_id=$7 AND code=$8 AND deleted_at IS NULL`,
		&model.Title, &model.CallbackURL, &model.Count, &model.Intervals, &model.Timeout, &model.Description,
		&model.UserId, &model.Code,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a setting in the database.
func (s settingRepository) Delete(userId int, code string) error {
	_, err := s.connection.Exec(
		s.context,
		"UPDATE ns_settings SET deleted_at=NOW() WHERE user_id=$1 AND code=$2 AND deleted_at IS NULL",
		&userId, &code,
	)

	return err
}

// FindByCode finds a setting in the database by code.
func (s settingRepository) FindByCode(userId int, code string) (*models.Setting, error) {
	var model models.Setting
	err := s.connection.QueryRow(
		s.context,
		`SELECT 
               id, code, user_id, COALESCE(resource_id, 0), title, callback_url, count, intervals, timeout, description, created_at 
             FROM 
               ns_settings 
             WHERE 
               user_id=$1 AND code=$2 AND deleted_at IS NULL`,
		userId, code,
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
func (s settingRepository) FindAll(userId int, limit int, offset int) (*map[int]models.Setting, error) {
	rows, err := s.connection.Query(
		s.context,
		`SELECT 
				id, code, user_id, COALESCE(resource_id, 0), title, callback_url, count, intervals, timeout, description, created_at
			FROM ns_settings s
			WHERE 
  				s.user_id=$1 AND s.deleted_at IS NULL 
			ORDER BY s.id ASC
			LIMIT $2 OFFSET $3`,
		userId, limit, offset,
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

// getModels returns array of the setting models.
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

// IsInDatabase setting is exists in the database.
func (s settingRepository) IsInDatabase(userId int, code string) (bool, error) {
	model, err := s.FindByCode(userId, code)

	return !(model == nil), err
}
