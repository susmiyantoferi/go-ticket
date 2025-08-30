package repository

import (
	"context"
	"errors"
	"ticket/domain/entity"
	"ticket/exception"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type EvenRepository interface {
	Create(ctx context.Context, event *entity.Event) (*entity.Event, error)
	Update(ctx context.Context, id uint, event *entity.Event) (*entity.Event, error)
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*entity.Event, error)
	FindAll(ctx context.Context, pg *entity.PaginateSearch) ([]*entity.Event, int64, error)
}

type evenRepositoryImpl struct {
	Db *gorm.DB
}

func NewEvenRepositoryImpl(db *gorm.DB) *evenRepositoryImpl {
	return &evenRepositoryImpl{
		Db: db,
	}
}

func (e *evenRepositoryImpl) Create(ctx context.Context, event *entity.Event) (*entity.Event, error) {
	if err := e.Db.WithContext(ctx).Create(event).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, exception.ErrorEventExist
		}
		return nil, err
	}

	return event, nil
}

func (e *evenRepositoryImpl) Update(ctx context.Context, id uint, event *entity.Event) (*entity.Event, error) {
	var events entity.Event
	if err := e.Db.WithContext(ctx).First(&events, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrorIdNotFound
		}
		return nil, err
	}

	if err := e.Db.WithContext(ctx).Model(&events).Updates(event).Error; err != nil {
		return nil, err
	}

	return &events, nil
}

func (e *evenRepositoryImpl) Delete(ctx context.Context, id uint) error {
	delete := e.Db.WithContext(ctx).Delete(&entity.Event{}, id)
	if delete.Error != nil {
		return delete.Error
	}

	if delete.RowsAffected == 0 {
		return exception.ErrorIdNotFound
	}

	return nil
}

func (e *evenRepositoryImpl) FindById(ctx context.Context, id uint) (*entity.Event, error) {
	var events entity.Event

	if err := e.Db.WithContext(ctx).First(&events, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrorIdNotFound
		}
		return nil, err
	}

	return &events, nil
}

func (e *evenRepositoryImpl) FindAll(ctx context.Context, pg *entity.PaginateSearch) ([]*entity.Event, int64, error) {
	var events []*entity.Event
	var totalItems int64

	query := e.Db.WithContext(ctx).Model(&entity.Event{})

	if pg.Search != "" {
		query = query.Where("name LIKE ?", "%"+pg.Search+"%")
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (pg.Page - 1) * pg.PageSize

	if err := query.Limit(pg.PageSize).Offset(offset).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, totalItems, nil
}
