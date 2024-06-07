package repository

import (
	"context"
	"fmt"
	"homework/internal/model"
	"homework/internal/pb"
)

type UserRepository struct {
	client pb.UserStorageClient
}

func NewUserRepository(client pb.UserStorageClient) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) AddUser(ctx context.Context, user model.User) error {
	req := &pb.AddUserRequest{User: &pb.User{
		Nickname:     user.Nickname,
		PasswordHash: user.PasswordHash,
	}}
	_, err := r.client.AddUser(ctx, req)
	if err != nil {
		return fmt.Errorf("grpc: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, nickname string) (model.User, error) {
	req := &pb.GetUserRequest{Nickname: nickname}
	res, err := r.client.GetUser(ctx, req)
	if err != nil {
		return model.User{}, fmt.Errorf("grpc: %w", err)
	}
	user := model.User{
		Nickname:     res.User.Nickname,
		PasswordHash: res.User.PasswordHash,
	}
	return user, nil
}
