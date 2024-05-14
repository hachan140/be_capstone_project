package services

import (
	"be-capstone-project/src/cmd/public/config"
	"be-capstone-project/src/pkg/adapter/repository/postgres"
	"be-capstone-project/src/pkg/adapter/repository/postgres/model"
	"be-capstone-project/src/pkg/core/common"
	"be-capstone-project/src/pkg/core/dtos/request"
	"be-capstone-project/src/pkg/core/dtos/response"
	"be-capstone-project/src/pkg/core/logger"
	"be-capstone-project/src/pkg/core/utils"
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
	accessToken, err := u.createToken(ttl, userModel.Email)
	if err != nil {
		return nil, err
	}
	res := &response.LoginResponse{AccessToken: accessToken}
	dat, err := u.validate(accessToken)
	if err != nil {
		return nil, err
	}
	fmt.Println(dat)
	return res, nil
}

func (u *UserService) createToken(ttl time.Duration, email string) (string, error) {
	privateKey := []byte(u.config.TokenConfig.AccessTokenPrivateKey)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}
	claims := &jwt.MapClaims{
		"iss":   "issuer",
		"email": email,
		"aud":   "audience",
		"exp":   time.Now().Add(ttl).Unix(),
		"iat":   time.Now().Unix(),
		"nbf":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Fatalf("Error signing token: %v", err)
	}
	return "Bearer " + tokenString, nil
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
	return claims["dat"], nil
}
