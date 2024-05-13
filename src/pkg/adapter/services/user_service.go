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
	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}
	var userModel *model.User
	var res *response.LoginResponse

	userModel, err = u.userRepository.FindUserByUsernameAndPassword(req.Email, hashedPassword)
	if err != nil {
		logger.Error(ctx, "Error when find user email", err)
		return nil, err
	}
	if userModel == nil {
		return nil, errors.New(common.ErrMessageInvalidUser)
	}
	accessToken, err := u.createToken(time.Duration(u.config.TokenConfig.AccessTokenTimeToLive), userModel.Email)
	if err != nil {
		return nil, err
	}
	res.AccessToken = accessToken
	return res, nil
}

func (u *UserService) createToken(ttl time.Duration, content interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(u.config.TokenConfig.AccessTokenPrivateKey))
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = content
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (u *UserService) validate(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(u.config.TokenConfig.AccessTokenPublicKey))
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims["dat"], nil
}
