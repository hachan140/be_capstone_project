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
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
	"time"
)

type IUserService interface {
	CreateUser(ctx context.Context, req *request.SignUpRequest) *common.ErrorCodeMessage
	LoginByUserEmail(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, *common.ErrorCodeMessage)
	LoginSocial(ctx context.Context, req *request.SocialLoginRequest) (*response.LoginResponse, *common.ErrorCodeMessage)
	RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.LoginResponse, *common.ErrorCodeMessage)
	ResetPasswordRequest(ctx context.Context, req *request.ResetPasswordRequest) (*string, *common.ErrorCodeMessage)
	ResetPassword(ctx context.Context, email string, req *request.ResetPassword) *common.ErrorCodeMessage
	UpdateUserStatusWhenEmailVerified(ctx context.Context, email string) *common.ErrorCodeMessage
}

type UserService struct {
	userRepository         postgres.IUserRepository
	refreshTokenRepository postgres.IRefreshTokenRepository
	config                 config.Config
}

func NewUserService(userRepository postgres.IUserRepository, refreshTokenRepo postgres.IRefreshTokenRepository, config config.Config) IUserService {
	return &UserService{userRepository: userRepository, refreshTokenRepository: refreshTokenRepo, config: config}
}

func (u *UserService) CreateUser(ctx context.Context, req *request.SignUpRequest) *common.ErrorCodeMessage {
	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	userEmailExisted, err := u.userRepository.FindUserByEmail(req.Email)
	if userEmailExisted != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeEmailExisted,
			Message:     common.ErrMessageEmailExisted,
		}
	}
	userModel := &model.User{
		FirstName: req.FistName,
		LastName:  req.LastName,
		Email:     req.Email,
		Status:    2,
		Password:  hashedPassword,
		Gender:    req.Gender,
		CreatedAt: time.Now(),
	}
	if err := u.userRepository.CreateUser(ctx, userModel); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if err := utils.SendVerifyEmailCreateAccount(u.config.VerifyEmailConfig.LinkVerifyEmail, u.config.EmailConfig.SenderEmail, u.config.EmailConfig.SenderPassword, userModel.Email); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (u *UserService) LoginByUserEmail(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, *common.ErrorCodeMessage) {
	var userModel *model.User
	userModel, err := u.userRepository.FindUserByEmail(req.Email)
	if err != nil {
		logger.Error(ctx, "Error when find user email", err)
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if userModel == nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeInvalidUser,
			Message:     common.ErrMessageInvalidUser,
		}
	}
	if userModel.Status != 1 {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserHaveNotVerifyEmail,
			Message:     common.ErrMessageUserHaveNotVerifyEmail,
		}
	}
	if err := utils.CheckHash(userModel.Password, req.Password); err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	ttl := time.Duration(u.config.TokenConfig.AccessTokenTimeToLive) * time.Second
	privateKey := u.config.TokenConfig.AccessTokenPrivateKey
	accessToken, err := u.createToken(ttl, userModel.Email, userModel.ID, "user_credentials", privateKey)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}

	refreshToken, err := u.generateRefreshToken(userModel.ID)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	res := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return res, nil
}

func (u *UserService) LoginSocial(ctx context.Context, req *request.SocialLoginRequest) (*response.LoginResponse, *common.ErrorCodeMessage) {
	var userModel *model.User
	userModel, err := u.userRepository.FindUserByEmail(req.Email)
	if err != nil {
		logger.Error(ctx, "Error when find user email", err)
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if userModel != nil && !userModel.IsSocial {
		if !userModel.IsSocial && userModel.Password == "" {
			return nil, &common.ErrorCodeMessage{
				HTTPCode:    http.StatusBadRequest,
				ServiceCode: common.ErrCodeUserSocialDoesNotExist,
				Message:     common.ErrMessageUserSocialDoesnotExist,
			}
		}
		if err := u.userRepository.UpdateUserSocial(userModel.ID); err != nil {
			return nil, &common.ErrorCodeMessage{
				HTTPCode:    http.StatusInternalServerError,
				ServiceCode: common.ErrCodeInternalError,
				Message:     err.Error(),
			}
		}
	}
	if userModel == nil {
		userModel = &model.User{
			FirstName:             req.FirstName,
			LastName:              req.LastName,
			Email:                 req.Email,
			Gender:                false,
			Status:                1,
			IsAdmin:               false,
			IsOrganizationManager: false,
			IsSocial:              req.IsSocial,
			CreatedAt:             time.Now(),
		}
		err := u.userRepository.CreateUser(ctx, userModel)
		if err != nil {
			return nil, &common.ErrorCodeMessage{
				HTTPCode:    http.StatusInternalServerError,
				ServiceCode: common.ErrCodeInternalError,
				Message:     err.Error(),
			}
		}
	}
	ttl := time.Duration(u.config.TokenConfig.AccessTokenTimeToLive) * time.Second
	privateKey := u.config.TokenConfig.AccessTokenPrivateKey
	accessToken, err := u.createToken(ttl, userModel.Email, userModel.ID, "socials", privateKey)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	refreshToken, err := u.generateRefreshToken(userModel.ID)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	res := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return res, nil
}

func (u *UserService) RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.LoginResponse, *common.ErrorCodeMessage) {
	rt, err := u.refreshTokenRepository.FindRefreshTokenByRefreshTokenString(req.RefreshToken)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if rt == nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeRefreshTokenNotFound,
			Message:     common.ErrMessageRefreshTokenNotFound,
		}
	}
	userModel, err := u.userRepository.FinduserByID(rt.UserID)
	if err != nil {
		logger.Error(ctx, "Error when find user email", err)
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if userModel == nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeInvalidUser,
			Message:     common.ErrMessageInvalidUser,
		}
	}
	ttl := time.Duration(u.config.TokenConfig.AccessTokenTimeToLive) * time.Second
	privateKey := u.config.TokenConfig.AccessTokenPrivateKey
	accessToken, err := u.createToken(ttl, userModel.Email, userModel.ID, "user_credentials", privateKey)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}

	refreshToken, err := u.generateRefreshToken(userModel.ID)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	res := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return res, nil
}

func (u *UserService) ResetPasswordRequest(ctx context.Context, req *request.ResetPasswordRequest) (*string, *common.ErrorCodeMessage) {
	var userModel *model.User
	userModel, err := u.userRepository.FindUserByEmail(req.Email)
	if err != nil {
		logger.Error(ctx, "Error when find user email", err)
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if userModel == nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeInvalidEmail,
			Message:     common.ErrMessageInvalidEmail,
		}
	}
	token, err := u.createToken(time.Duration(u.config.ResetPasswordConfig.ResetPasswordTokenTTL)*time.Second, req.Email, userModel.ID, "reset_password", u.config.ResetPasswordConfig.ResetPasswordPrivateKey)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	resetPasswordLink := fmt.Sprintf("%s?token=%s", "http://localhost:5173/reset-password", token)
	err = utils.SendResetPasswordLink(resetPasswordLink, u.config.ResetPasswordConfig.ResetPasswordSender, u.config.ResetPasswordConfig.ResetPasswordSenderPassword, req.Email)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil, nil
}

func (u *UserService) ResetPassword(ctx context.Context, email string, req *request.ResetPassword) *common.ErrorCodeMessage {
	var userModel *model.User
	userModel, err := u.userRepository.FindUserByEmail(email)
	if err != nil {
		logger.Error(ctx, "Error when find user email", err)
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if userModel == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeInvalidEmail,
			Message:     common.ErrMessageInvalidEmail,
		}
	}
	hashedPassword, err := utils.EncryptPassword(req.NewPassword)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if err := u.userRepository.ResetPassword(userModel.ID, hashedPassword); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (u *UserService) UpdateUserStatusWhenEmailVerified(ctx context.Context, email string) *common.ErrorCodeMessage {
	var userModel *model.User
	userModel, err := u.userRepository.FindUserByEmail(email)
	if err != nil {
		logger.Error(ctx, "Error when find user email", err)
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if userModel == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeInvalidEmail,
			Message:     common.ErrMessageInvalidEmail,
		}
	}
	if err := u.userRepository.UpdateUserStatus(userModel.ID, 1); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (u *UserService) createToken(ttl time.Duration, email string, userID uint, aud string, prvKey string) (string, error) {
	privateKey := []byte(prvKey)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}
	claims := &jwt.MapClaims{
		"iss":     "GeniFast-Search_Go",
		"email":   email,
		"user_id": fmt.Sprintf("%v", userID),
		"aud":     aud,
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

func (u *UserService) generateRefreshToken(userID uint) (string, error) {
	rt, err := u.refreshTokenRepository.FindRefreshTokenByUserID(userID)
	if err != nil {
		return "", err
	}
	refreshTokenString := utils.RandomString(128)
	if rt == nil {
		rt = &model.RefreshToken{
			UserID:       userID,
			RefreshToken: refreshTokenString,
		}
		err := u.refreshTokenRepository.CreateRefreshToken(rt)
		if err != nil {
			return "", err
		}
	} else {
		if err := u.refreshTokenRepository.UpdateRefreshToken(rt.ID, refreshTokenString); err != nil {
			return "", err
		}
	}
	return refreshTokenString, nil
}
