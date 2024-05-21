package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/storage"
)

type IRefreshTokenRepository interface {
	CreateRefreshToken(rt *model.RefreshToken) error
	UpdateRefreshToken(rtID uint, refreshToken string) error
	FindRefreshTokenByUserID(userID uint) (*model.RefreshToken, error)
	FindRefreshTokenByRefreshTokenString(refreshToken string) (*model.RefreshToken, error)
}

type RefreshTokenRepository struct {
	storage *storage.Database
}

func NewRefreshTokenRepository(storage *storage.Database) IRefreshTokenRepository {
	return &RefreshTokenRepository{storage: storage}
}

func (r *RefreshTokenRepository) CreateRefreshToken(rt *model.RefreshToken) error {
	err := r.storage.Model(rt).Create(&rt).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) UpdateRefreshToken(rtID uint, refreshToken string) error {
	err := r.storage.Exec("update refresh_tokens set refresh_token = ? where id = ?", refreshToken, rtID).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) FindRefreshTokenByUserID(userID uint) (*model.RefreshToken, error) {
	var rt *model.RefreshToken
	err := r.storage.Raw("select * from refresh_tokens where user_id = ?", userID).Scan(&rt).Error
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (r *RefreshTokenRepository) FindRefreshTokenByRefreshTokenString(refreshToken string) (*model.RefreshToken, error) {
	var rt *model.RefreshToken
	err := r.storage.Raw("select * from refresh_tokens where refresh_token = ?", refreshToken).Scan(&rt).Error
	if err != nil {
		return nil, err
	}
	return rt, nil
}
