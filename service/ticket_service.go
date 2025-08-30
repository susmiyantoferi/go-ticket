package service

import (
	"context"
	"errors"
	"fmt"
	"ticket/domain/entity"
	"ticket/exception"
	"ticket/repository"

	"github.com/go-playground/validator/v10"
)

type TicketService interface {
	Create(ctx context.Context, req *entity.TicketCreateRequest) (*entity.TicketResponse, error)
	Update(ctx context.Context, id uint, req *entity.TicketUpdateRequest) (*entity.TicketResponse, error)
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*entity.TicketResponse, error)
	FindByUserId(ctx context.Context, userId uint) ([]*entity.TicketResponse, error)
	FindAll(ctx context.Context) ([]*entity.TicketResponse, error)
}

type ticketServiceImpl struct {
	TicketRepo repository.TicketRepository
	EventRepo  repository.EvenRepository
	Validate   *validator.Validate
}

func NewTicketServiceImpl(ticketRepo repository.TicketRepository, eventRepo repository.EvenRepository, validate *validator.Validate) *ticketServiceImpl {
	return &ticketServiceImpl{
		TicketRepo: ticketRepo,
		EventRepo:  eventRepo,
		Validate:   validate,
	}
}

const (
	Confirm string = "confirm"
	Cancel  string = "cancel"
	Waiting string = "waiting"
)

func (t *ticketServiceImpl) Create(ctx context.Context, req *entity.TicketCreateRequest) (*entity.TicketResponse, error) {
	if err := t.Validate.Struct(req); err != nil {
		return nil, exception.ErrorValidation
	}

	event, err := t.EventRepo.FindById(ctx, req.EventID)
	if err != nil {
		if errors.Is(err, exception.ErrorIdNotFound) {
			return nil, exception.ErrorEventNotFound
		}
		return nil, fmt.Errorf("ticket service: find event: %w", err)
	}

	if event.Capacity < req.Qty || req.Qty < 0 {
		if errors.Is(err, exception.ErrorQty) {
			return nil, exception.ErrorQty
		}
		return nil, exception.ErrorStockNotEnough
	}

	total := float64(req.Qty) * event.Price

	ticks := entity.Ticket{
		UserID:      req.UserID,
		EventID:     event.ID,
		Qty:         req.Qty,
		UnitPrice:   event.Price,
		TotalAmount: total,
		Status:      Waiting,
	}

	result, err := t.TicketRepo.Create(ctx, &ticks, event)
	if err != nil {
		return nil, fmt.Errorf("ticket service: create: %w", err)
	}

	response := entity.ToTicketResponse(result)

	return response, nil
}

func (t *ticketServiceImpl) Update(ctx context.Context, id uint, req *entity.TicketUpdateRequest) (*entity.TicketResponse, error) {
	if err := t.Validate.Struct(req); err != nil {
		return nil, exception.ErrorValidation
	}

	ticks := entity.Ticket{
		Status: req.Status,
	}

	result, err := t.TicketRepo.Update(ctx, id, &ticks)
	if err != nil {
		if errors.Is(err, exception.ErrorIdNotFound) {
			return nil, exception.ErrorIdNotFound
		}
		return nil, fmt.Errorf("ticket service: update: %w", err)
	}

	response := entity.ToTicketResponse(result)

	return response, nil
}

func (t *ticketServiceImpl) Delete(ctx context.Context, id uint) error {
	if err := t.TicketRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, exception.ErrorIdNotFound) {
			return exception.ErrorIdNotFound
		}
		return fmt.Errorf("ticket service: delete: %w", err)
	}

	return nil
}

func (t *ticketServiceImpl) FindById(ctx context.Context, id uint) (*entity.TicketResponse, error) {
	result, err := t.TicketRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, exception.ErrorIdNotFound) {
			return nil, exception.ErrorIdNotFound
		}
		return nil, fmt.Errorf("ticket service: find by id: %w", err)
	}

	response := entity.ToTicketResponse(result)

	return response, nil
}

func (t *ticketServiceImpl) FindByUserId(ctx context.Context, userId uint) ([]*entity.TicketResponse, error) {
	result, err := t.TicketRepo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("ticket service: find by user id: %w", err)
	}

	var responses []*entity.TicketResponse
	for _, v := range result {
		response := entity.TicketResponse{
			UserID: v.UserID,
			User: entity.UserInfo{
				Email:   v.User.Email,
				Name:    v.User.Name,
				Hp:      v.User.Hp,
				Address: v.User.Address,
			},
			EventID: v.EventID,
			Event: entity.EventInfo{
				Name:        v.Event.Name,
				Description: v.Event.Description,
			},
			Qty:         v.Qty,
			UnitPrice:   v.UnitPrice,
			TotalAmount: v.TotalAmount,
			Status:      v.Status,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}

func (t *ticketServiceImpl) FindAll(ctx context.Context) ([]*entity.TicketResponse, error) {
	result, err := t.TicketRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("ticket service: find all: %w", err)
	}

	var responses []*entity.TicketResponse
	for _, v := range result {
		response := entity.TicketResponse{
			UserID: v.UserID,
			User: entity.UserInfo{
				Email:   v.User.Email,
				Name:    v.User.Name,
				Hp:      v.User.Hp,
				Address: v.User.Address,
			},
			EventID: v.EventID,
			Event: entity.EventInfo{
				Name:        v.Event.Name,
				Description: v.Event.Description,
			},
			Qty:         v.Qty,
			UnitPrice:   v.UnitPrice,
			TotalAmount: v.TotalAmount,
			Status:      v.Status,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}
