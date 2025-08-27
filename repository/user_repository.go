package repository

import (
	"context"
	"errors"
	"ticket/domain/entity"
	"ticket/exception"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserReposiitory interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, userId uint, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, userId uint) error
	FindById(ctx context.Context, userId uint) (*entity.User, error)
	FindAll(ctx context.Context) ([]*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userReposiitoryImpl struct {
	Db *gorm.DB
}

func NewUserReposiitoryImpl(db *gorm.DB) *userReposiitoryImpl {
	return &userReposiitoryImpl{
		Db: db,
	}
}

func (u *userReposiitoryImpl) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := u.Db.WithContext(ctx).Create(user).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, exception.ErrorEmailExist
		}
		return nil, err
	}

	return user, nil
}

func (u *userReposiitoryImpl) Update(ctx context.Context, userId uint, user *entity.User) (*entity.User, error) {
	var users entity.User
	if err := u.Db.WithContext(ctx).First(&users, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrorIdNotFound
		}

		return nil, err
	}

	update := map[string]interface{}{}

	if user.Name != "" {
		update["name"] = user.Name
	}

	if user.Password != "" {
		update["password"] = user.Password
	}

	if user.Hp != "" {
		update["hp"] = user.Hp
	}

	if user.Address != "" {
		update["address"] = user.Address
	}

	if err := u.Db.WithContext(ctx).Model(&users).Updates(update).Error; err != nil {
		return nil, err
	}

	return &users, nil
}

func (u *userReposiitoryImpl) Delete(ctx context.Context, userId uint) error {
	result := u.Db.WithContext(ctx).Delete(&entity.User{}, userId)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exception.ErrorIdNotFound
	}

	return nil
}

func (u *userReposiitoryImpl) FindById(ctx context.Context, userId uint) (*entity.User, error) {
	var user entity.User
	if err := u.Db.WithContext(ctx).First(&user, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrorIdNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (u *userReposiitoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := u.Db.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrorEmailNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (u *userReposiitoryImpl) FindAll(ctx context.Context) ([]*entity.User, error) {
	var user []*entity.User

	result := u.Db.WithContext(ctx).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
