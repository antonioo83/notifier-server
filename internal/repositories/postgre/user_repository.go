package postgre

import (
	"context"
	"errors"
	"github.com/antonioo83/notifier-server/internal/models"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewUserRepository(context context.Context, pool *pgxpool.Pool) interfaces.UserRepository {
	return &userRepository{context, pool}
}

func (u userRepository) Save(user models.User) error {
	var lastInsertId int
	err := u.connection.QueryRow(
		u.context,
		"INSERT INTO ns_users(code, role, title, auth_token, description)VALUES ($1, $2, $3, $4, $5) RETURNING id",
		&user.Code, &user.Role, &user.Title, &user.AuthToken, &user.Description,
	).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	return err
}

func (u userRepository) Update(model models.User) error {
	_, err := u.connection.Exec(
		u.context,
		"UPDATE ns_users SET role=$1, title=$2, auth_token=$3, description=$4, updated_at=NOW() WHERE code=$5 AND deleted_at IS NULL",
		&model.Role, &model.Title, &model.AuthToken, &model.Description, &model.Code,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u userRepository) Delete(code string) error {
	_, err := u.connection.Exec(u.context, "UPDATE ns_users SET deleted_at=NOW() WHERE code=$1 AND deleted_at IS NULL", code)

	return err
}

func (u userRepository) FindByCode(code string) (*models.User, error) {
	var model models.User
	err := u.connection.QueryRow(
		u.context,
		"SELECT id, code, role, title, auth_token, description, created_at FROM ns_users WHERE code=$1 AND deleted_at IS NULL",
		code,
	).Scan(&model.ID, &model.Code, &model.Role, &model.Title, &model.AuthToken, &model.Description, &model.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (u userRepository) FindByToken(code string) (*models.User, error) {
	var model models.User
	err := u.connection.QueryRow(
		u.context,
		"SELECT id, code, role, title, auth_token, description, created_at FROM ns_users WHERE auth_token=$1 AND deleted_at IS NULL",
		code,
	).Scan(&model.ID, &model.Code, &model.Role, &model.Title, &model.AuthToken, &model.Description, &model.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (u userRepository) FindAll(limit int, offset int) (*map[int]models.User, error) {
	rows, err := u.connection.Query(
		u.context,
		`SELECT 
				u.id, u.code, u.role, u.title, u.description
			FROM ns_users u
			WHERE 
  				u.deleted_at IS NULL 
			ORDER BY u.id ASC
			LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}

	users, err := getModels(rows)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func getModels(rows pgx.Rows) (map[int]models.User, error) {
	var users = make(map[int]models.User)
	var model models.User
	var user = models.User{}
	lastUserId := 0
	for rows.Next() {
		err := rows.Scan(
			&model.ID, &model.Code, &model.Role, &model.Title, &model.Description, //&model.CreatedAt, &model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if lastUserId != model.ID {
			user = models.User{}
			user.ID = model.ID
			user.Role = model.Role
			user.Code = model.Code
			user.Title = model.Title
			user.Description = model.Description
		}
		lastUserId = model.ID
		users[user.ID] = user
	}

	return users, nil
}

func (u userRepository) IsInDatabase(code string) (bool, error) {
	model, err := u.FindByCode(code)

	return !(model == nil), err
}
