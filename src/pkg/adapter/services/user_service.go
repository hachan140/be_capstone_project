package services

import (
	"be-capstone-project/src/pkg/adapter/repository/postgres"
	"be-capstone-project/src/pkg/adapter/repository/postgres/model"
	"be-capstone-project/src/pkg/core/dtos/request"
	"be-capstone-project/src/pkg/core/utils"
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
	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return err
	}
	userModel := &model.User{
		FirstName: req.FistName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
		Gender:    req.Gender,
		CreatedAt: time.Now(),
	}
	if err := u.userRepository.CreateUser(ctx, userModel); err != nil {
		return err
	}
	return nil
}
