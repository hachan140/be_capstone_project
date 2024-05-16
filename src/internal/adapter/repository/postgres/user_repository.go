package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/storage"
	"context"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	FindUserByEmailAndPassword(username string, password string) (*model.User, error)
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
	result := u.storage.Raw("select * from users u where u.email = ?", email).Scan(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (u *UserRepository) FindUserByEmailAndPassword(email string, password string) (*model.User, error) {
	var user *model.User
	result := u.storage.Raw("select * from users where email = ? and password = ?", email, password).Scan(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
