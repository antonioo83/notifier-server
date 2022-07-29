package postgre

import (
	"context"
	"errors"
	"github.com/antonioo83/notifier-server/internal/models"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type messageRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewMessageRepository(context context.Context, pool *pgxpool.Pool) interfaces.MessageRepository {
	return &messageRepository{context, pool}
}

func (u messageRepository) Save(model models.Message) error {
	var lastInsertId int
	err := u.connection.QueryRow(
		u.context,
		`INSERT INTO ns_messages(code, user_id, resource_id, command, priority, content, is_sent, attempt_count, 
				is_sent_callback, callback_attempt_count, description, send_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
		&model.Code, &model.UserId, &model.ResourceId, &model.Command, &model.Priority, &model.Content, &model.IsSent,
		&model.AttemptCount, &model.IsSentCallback, &model.CallbackAttemptCount, &model.Description, &model.SendAt,
	).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	return err
}

func (u messageRepository) Update(model models.Message) error {
	_, err := u.connection.Exec(
		u.context,
		`UPDATE 
               ns_messages 
             SET 
               user_id=$1, 
               resource_id=$2, 
               command=$3, 
               priority=$4, 
               content=$5, 
               is_sent=$6, 
               attempt_count=$7, 
               is_sent_callback=$8, 
               callback_attempt_count=$9, 
               description=$10, 
               send_at=$11, 
               updated_at=NOW() 
             WHERE 
               code=$12 AND deleted_at IS NULL`,
		&model.UserId, &model.ResourceId, &model.Command, &model.Priority, &model.Content, &model.IsSent,
		&model.AttemptCount, &model.IsSentCallback, &model.CallbackAttemptCount, &model.Description, &model.SendAt,
		&model.Code,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u messageRepository) Delete(code string) error {
	_, err := u.connection.Exec(u.context, "UPDATE ns_messages SET deleted_at=NOW() WHERE code=$1 AND deleted_at IS NULL", code)

	return err
}

func (u messageRepository) FindByCode(code string) (*models.Message, error) {
	var model models.Message
	err := u.connection.QueryRow(
		u.context,
		`SELECT r.url, m.id, m.code, m.user_id, m.resource_id, m.command, m.priority, m.content, m.is_sent, m.attempt_count, 
				m.is_sent_callback, m.callback_attempt_count, m.description, m.send_at, m.created_at 
             FROM 
               ns_messages m
             LEFT JOIN ns_resources r ON r.id=m.resource_id 
             WHERE 
               code=$1 AND deleted_at IS NULL`,
		code,
	).Scan(&model.Resource.URL, &model.ID, &model.Code, &model.UserId, &model.ResourceId, &model.Command, &model.Priority,
		&model.Content, &model.IsSent, &model.AttemptCount, &model.IsSentCallback, &model.CallbackAttemptCount,
		&model.Description, &model.SendAt, &model.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (u messageRepository) FindAll(limit int, offset int) (*map[int]models.Message, error) {
	rows, err := u.connection.Query(
		u.context,
		`SELECT r.url, m.id, m.code, m.user_id, m.resource_id, m.command, m.priority, m.content, m.is_sent, m.attempt_count, 
				m.is_sent_callback, m.callback_attempt_count, m.description, m.send_at, m.created_at 
             FROM 
               ns_messages m
             LEFT JOIN ns_resources r ON r.id=m.resource_id 
             WHERE 
  				m.deleted_at IS NULL 
			ORDER BY m.id ASC
			LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}

	users, err := getMessageModels(rows)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func getMessageModels(rows pgx.Rows) (map[int]models.Message, error) {
	var messages = make(map[int]models.Message)
	for rows.Next() {
		var model models.Message
		model.Resource = models.Resource{}
		err := rows.Scan(
			&model.Resource.URL, &model.ID, &model.Code, &model.UserId, &model.ResourceId, &model.Command, &model.Priority,
			&model.Content, &model.IsSent, &model.AttemptCount, &model.IsSentCallback, &model.CallbackAttemptCount,
			&model.Description, &model.SendAt, &model.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages[model.ID] = model
	}

	return messages, nil
}

func (u messageRepository) IsInDatabase(code string) (bool, error) {
	model, err := u.FindByCode(code)

	return !(model == nil), err
}
