package services

import (
	"be-capstone-project/src/cmd/public/config"
	"be-capstone-project/src/internal/adapter/repository/postgres"
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/dtos/response"
	"be-capstone-project/src/internal/core/logger"
	"be-capstone-project/src/internal/core/utils"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"strings"
	"time"
)

type IUserService interface {
	CreateUser(ctx context.Context, req *request.SignUpRequest) error
	LoginByUserEmail(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error)
}

type UserService struct {
	userRepository postgres.IUserRepository
	config         config.Config
}

func NewUserService(userRepository postgres.IUserRepository, config config.Config) IUserService {
	return &UserService{userRepository: userRepository, config: config}
}

func (u *UserService) CreateUser(ctx context.Context, req *request.SignUpRequest) error {
	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return err
	}
	userEmailExisted, err := u.userRepository.FindUserByEmail(req.Email)
	if userEmailExisted != nil {
		return errors.New(common.ErrMessageEmailExisted)
	}
	userModel := &model.User{
		FirstName: req.FistName,
		LastName:  req.LastName,
		Email:     req.Email,
		Status:    1,
		Password:  hashedPassword,
		Gender:    req.Gender,
		CreatedAt: time.Now(),
	}
	if err := u.userRepository.CreateUser(ctx, userModel); err != nil {
		return err
	}
	return nil
}

func (u *UserService) LoginByUserEmail(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	var userModel *model.User
	userModel, err := u.userRepository.FindUserByEmail(req.Email)
	if err != nil {
		logger.Error(ctx, "Error when find user email", err)
		return nil, err
	}
	if userModel == nil {
		return nil, errors.New(common.ErrMessageInvalidUser)
	}
	if err := utils.CheckHash(userModel.Password, req.Password); err != nil {
		return nil, errors.New(common.ErrMessageInvalidPassword)
	}
	ttl := time.Duration(u.config.TokenConfig.AccessTokenTimeToLive) * time.Second
	accessToken, err := u.createToken(ttl, userModel.Email, userModel.ID)
	if err != nil {
		return nil, err
	}

	res := &response.LoginResponse{AccessToken: accessToken}
	return res, nil
}

func (u *UserService) createToken(ttl time.Duration, email string, userID uint) (string, error) {
	privateKey := []byte(u.config.TokenConfig.AccessTokenPrivateKey)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}
	claims := &jwt.MapClaims{
		"iss":     "GeniFast-Search_Go",
		"email":   email,
		"user_id": fmt.Sprintf("%v", userID),
		"aud":     "user_credentials",
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Fatalf("Error signing token: %v", err)
	}
	return tokenString, nil
}

func (u *UserService) validate(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(u.config.TokenConfig.AccessTokenPublicKey))
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}
	args := strings.Split(token, " ")
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(args[1], claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	return claims["email"], nil
}
