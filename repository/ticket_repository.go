package repository

import (
	"context"
	"errors"
	"ticket/domain/entity"
	"ticket/exception"

	"gorm.io/gorm"
)

type TicketRepository interface {
	Create(ctx context.Context, ticket *entity.Ticket, ev *entity.Event) (*entity.Ticket, error)
	Update(ctx context.Context, id uint, ticket *entity.Ticket) (*entity.Ticket, error)
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*entity.Ticket, error)
	FindByUserId(ctx context.Context, userId uint) ([]*entity.Ticket, error)
	FindAll(ctx context.Context) ([]*entity.Ticket, error)
}

type ticketRepositoryImpl struct {
	Db *gorm.DB
}

func NewTicketRepositoryImpl(db *gorm.DB) *ticketRepositoryImpl {
	return &ticketRepositoryImpl{
		Db: db,
	}
}

func (t *ticketRepositoryImpl) Create(ctx context.Context, ticket *entity.Ticket, ev *entity.Event) (*entity.Ticket, error) {

	err := t.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.WithContext(ctx).Create(ticket).Error; err != nil {
			return err
		}

		reduceCaps := tx.WithContext(ctx).Model(ev).Where("id = ? ", ev.ID).
			UpdateColumn("capacity", gorm.Expr("capacity - ?", ticket.Qty))
		if reduceCaps.Error != nil {
			return reduceCaps.Error
		}

		if reduceCaps.RowsAffected == 0 {
			return exception.ErrorQty
		}

		if err := tx.WithContext(ctx).Preload("User").Preload("Event").First(&ticket, ticket.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (t *ticketRepositoryImpl) Update(ctx context.Context, id uint, ticket *entity.Ticket) (*entity.Ticket, error) {
	var ticks entity.Ticket
	if err := t.Db.WithContext(ctx).First(&ticks, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrorIdNotFound
		}
	}

	if err := t.Db.WithContext(ctx).Model(&ticks).Updates(ticket).Error; err != nil {
		return nil, err
	}

	if err := t.Db.WithContext(ctx).Preload("User").Preload("Event").First(&ticks, ticks.ID).Error; err != nil {
		return nil, err
	}

	return &ticks, nil
}

func (t *ticketRepositoryImpl) Delete(ctx context.Context, id uint) error {
	delete := t.Db.WithContext(ctx).Delete(&entity.Ticket{}, id)
	if delete.Error != nil {
		return delete.Error
	}

	if delete.RowsAffected == 0 {
		return exception.ErrorIdNotFound
	}

	return nil
}

func (t *ticketRepositoryImpl) FindById(ctx context.Context, id uint) (*entity.Ticket, error) {
	var ticks entity.Ticket
	if err := t.Db.WithContext(ctx).Preload("User").Preload("Event").First(&ticks, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrorIdNotFound
		}
		return nil, err
	}

	return &ticks, nil
}

func (t *ticketRepositoryImpl) FindByUserId(ctx context.Context, userId uint) ([]*entity.Ticket, error) {
	var ticks []*entity.Ticket

	result := t.Db.WithContext(ctx).Preload("User").Preload("Event").Where("user_id = ?", userId).Find(&ticks)
	if result.Error != nil {
		return nil, result.Error
	}

	return ticks, nil

}

func (t *ticketRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Ticket, error) {
	var ticks []*entity.Ticket
	if err := t.Db.WithContext(ctx).Preload("User").Preload("Event").Find(&ticks).Error; err != nil {
		return nil, err
	}

	return ticks, nil
}
