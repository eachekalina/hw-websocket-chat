package repository

import (
	"context"
	"pg-course/internal/domain"
	"pg-course/internal/repository/queries"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type repo struct {
	*queries.Queries
	pool   *pgxpool.Pool
	logger logrus.FieldLogger
}

func NewRepository(pgxPool *pgxpool.Pool, logger logrus.FieldLogger) UserRepository {
	return &repo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
		logger:  logger,
	}
}

type UserRepository interface {
	FindUserByID(ctx context.Context, id int) (*domain.User, error)
}
