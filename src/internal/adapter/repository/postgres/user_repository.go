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
	FinduserByID(userID uint) (*model.User, error)
	UpdateUserOrganizationRole(userID uint, orgID uint, isOrgManager bool) error
	UpdateUserSocial(userID uint) error
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

func (u *UserRepository) FinduserByID(userID uint) (*model.User, error) {
	var user *model.User
	result := u.storage.Raw("select * from users u where u.id = ?", userID).Scan(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (u *UserRepository) UpdateUserOrganizationRole(userID uint, orgID uint, isOrgManager bool) error {
	err := u.storage.Exec("update users set is_organization_manager = ?, organization_id = ? where id = ?", isOrgManager, orgID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) UpdateUserSocial(userID uint) error {
	err := u.storage.Exec("update users set is_social = true where id = ?", userID).Error
	if err != nil {
		return err
	}
	return nil
}
