package user

import (
	"github.com/jinzhu/gorm"
	"skripsi-sso/database/entities"
)

type userRepo struct {
	db []*gorm.DB
}

func NewUserRepo(db []*gorm.DB) *userRepo {
	return &userRepo{db: db}
}

func (r userRepo) SaveUser(entity *entities.User) error {
	err := r.db[0].Create(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r userRepo) UpdateUser(entity *entities.User) error {
	err := r.db[0].Save(&entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r userRepo) FindUser(entity *entities.User, offset int, limit int) *[]entities.User {
	var users []entities.User
	err := r.db[0].Find(&users).Offset(offset).Limit(limit).Error
	if err != nil {
		return nil
	}

	return &users
}

func (r userRepo) FindUserByEmailForAuthSso(email string) (*entities.UserAuth, error) {
	var user entities.User
	err := r.db[0].
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	userAuth := entities.UserAuth{}
	userAuth.Email = user.Email
	userAuth.Name = user.Name
	userAuth.Phone = user.Phone
	userAuth.Address = user.Address
	userAuth.Password = user.Password

	return &userAuth, nil
}

func (r userRepo) FindUserById(entity *entities.User) (*entities.User, error) {
	var user entities.User
	err := r.db[0].Where("user_id = ?", entity.Id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r userRepo) DeleteUser(entity *entities.User) error {
	err := r.db[0].Delete(&entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r userRepo) CountUser(entity entities.User) (int64, error) {
	var count int64
	err := r.db[0].Table(entity.TableName()).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
