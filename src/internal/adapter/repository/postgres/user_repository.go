package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/storage"
	"context"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
}

type UserRepository struct {
	storage *storage.Database
}

func NewUserRepository(storage *storage.Database) IUserRepository {
	return &UserRepository{storage: storage}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	result := u.storage.Model(&user).Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
