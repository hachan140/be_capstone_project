package services

import (
	"be-capstone-project/src/internal/adapter/repository/postgres"
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/dtos/request"
	"context"
	"time"
)

type IUserService interface {
	CreateUser(ctx context.Context, req *request.SignUpRequest) error
}

type UserService struct {
	userRepository postgres.IUserRepository
}

func NewUserService(userRepository postgres.IUserRepository) IUserService {
	return &UserService{userRepository: userRepository}
}

func (u *UserService) CreateUser(ctx context.Context, req *request.SignUpRequest) error {
	// TODO: hash password
	userModel := &model.User{
		FirstName: req.FistName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Gender:    req.Gender,
		CreatedAt: time.Now(),
	}
	err := u.userRepository.CreateUser(ctx, userModel)
	if err != nil {
		return err
	}
	return nil
}
