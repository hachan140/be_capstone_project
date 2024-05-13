package postgres

import (
	"be-capstone-project/src/pkg/adapter/repository/postgres/model"
	"be-capstone-project/src/pkg/core/storage"
	"context"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	FindUserByUsernameAndPassword(username string, password string) (*model.User, error)
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

func (u *UserRepository) FindUserByEmail(email string) (*model.User, error) {
	var user *model.User
	result := u.storage.Where("email = ?", email).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (u *UserRepository) FindUserByUsernameAndPassword(username string, password string) (*model.User, error) {
	var user *model.User
	result := u.storage.Where("username = ? and password", username, password).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
