package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/model"
)

type MessageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{pool: pool}
}

func (r *MessageRepository) AddMessage(ctx context.Context, message model.Message) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO messages (sent_time, nickname, message) VALUES ($1, $2, $3)`, message.SentTime, message.Nickname, message.Message)
	return err
}

func (r *MessageRepository) GetLastMessages(ctx context.Context, n int) ([]model.Message, error) {
	rows, err := r.pool.Query(ctx, `SELECT sent_time, nickname, message FROM messages ORDER BY sent_time DESC LIMIT $1`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Message])
}