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
	FindUsersInOrganization(emails []*string) ([]string, error)
	FindUsersNotInOrganization(emails []*string) ([]*model.User, error)
	AddPeopleOrganization(userID uint, orgID uint, deptID uint) error
	ResetPassword(userID uint, password string) error
	UpdateUserStatus(userID uint, status int) error
	FindUserInOrganization(email string, orgID uint) (*model.User, error)
	UpdateUserRoleManager(userID uint, isManager bool) error
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

func (u *UserRepository) FindUsersInOrganization(emails []*string) ([]string, error) {
	var emailExisted []string
	result := u.storage.Raw("select email from users u where email in ? and organization_id != 0", emails).Scan(&emailExisted)
	if result.Error != nil {
		return nil, result.Error
	}
	return emailExisted, nil
}

func (u *UserRepository) FindUsersNotInOrganization(emails []*string) ([]*model.User, error) {
	var users []*model.User
	result := u.storage.Raw("select * from users u where email in ? and organization_id = 0 and status = 1", emails).Scan(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil

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

func (u *UserRepository) AddPeopleOrganization(userID uint, orgID uint, deptID uint) error {
	err := u.storage.Exec("update users set organization_id = ?, dept_id = ? where id = ?", orgID, deptID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) ResetPassword(userID uint, password string) error {
	err := u.storage.Exec("update users set password = ? where id = ?", password, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) UpdateUserStatus(userID uint, status int) error {
	err := u.storage.Exec("update users set status = ? where id = ?", status, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) FindUserInOrganization(email string, orgID uint) (*model.User, error) {
	var user *model.User
	result := u.storage.Raw("select * from users u where email = ? and organization_id = ?", email, orgID).Scan(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil

}

func (u *UserRepository) UpdateUserRoleManager(userID uint, isManager bool) error {
	err := u.storage.Exec("update users set is_organization_manager = ? where id = ?", isManager, userID).Error
	if err != nil {
		return err
	}
	return nil
}
