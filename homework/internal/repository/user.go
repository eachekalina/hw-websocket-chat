package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/model"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) AddUser(ctx context.Context, user model.User) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO users (nickname, password_hash) VALUES ($1, $2)`, user.Nickname, user.PasswordHash)
	return err
}

func (r *UserRepository) GetUser(ctx context.Context, nickname string) (model.User, error) {
	row := r.pool.QueryRow(ctx, `SELECT password_hash FROM users WHERE nickname = $1`, nickname)

	var passwordHash []byte
	err := row.Scan(&passwordHash)
	if err != nil {
		return model.User{}, err
	}

	return model.User{Nickname: nickname, PasswordHash: passwordHash}, nil
}
