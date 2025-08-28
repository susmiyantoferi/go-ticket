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

type EventService interface {
	Create(ctx context.Context, req *entity.EventCreateRequest) (*entity.EventResponse, error)
	Update(ctx context.Context, id uint, req *entity.EventUpdateRequest) (*entity.EventResponse, error)
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*entity.EventResponse, error)
	FindAll(ctx context.Context) ([]*entity.EventResponse, error)
}

type eventServiceImpl struct {
	EventRepo repository.EvenRepository
	Validate  *validator.Validate
}

func NewEventServiceImpl(eventRepo repository.EvenRepository, validate *validator.Validate) *eventServiceImpl {
	return &eventServiceImpl{
		EventRepo: eventRepo,
		Validate:  validate,
	}
}

const (
	Aktif       string = "aktif"
	Berlangsung string = "berlangsung"
	Selesai     string = "selesai"
)

func (e *eventServiceImpl) Create(ctx context.Context, req *entity.EventCreateRequest) (*entity.EventResponse, error) {
	if err := e.Validate.Struct(req); err != nil {
		return nil, exception.ErrorValidation
	}

	events := entity.Event{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Capacity:    req.Capacity,
		Status:      Aktif,
	}

	result, err := e.EventRepo.Create(ctx, &events)
	if err != nil {
		if errors.Is(err, exception.ErrorEventExist) {
			return nil, exception.ErrorEventExist
		}
		return nil, fmt.Errorf("event service: create: %w", err)
	}

	response := entity.ToEventResponse(result)

	return response, nil
}

func (e *eventServiceImpl) Update(ctx context.Context, id uint, req *entity.EventUpdateRequest) (*entity.EventResponse, error) {
	if err := e.Validate.Struct(req); err != nil {
		return nil, exception.ErrorValidation
	}

	events := entity.Event{}

	if req.Name != nil {
		events.Name = *req.Name
	}

	if req.Description != nil {
		events.Description = *req.Description
	}

	if req.Capacity != nil {
		events.Capacity = *req.Capacity
	}

	if req.Price != nil {
		events.Price = *req.Price
	}

	if req.Status != nil {
		events.Status = *req.Status
	}

	result, err := e.EventRepo.Update(ctx, id, &events)
	if err != nil {
		if errors.Is(err, exception.ErrorIdNotFound) {
			return nil, exception.ErrorIdNotFound
		}
		return nil, fmt.Errorf("event service: update: %w", err)
	}

	response := entity.ToEventResponse(result)

	return response, nil
}

func (e *eventServiceImpl) Delete(ctx context.Context, id uint) error {
	if err := e.EventRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, exception.ErrorIdNotFound) {
			return exception.ErrorIdNotFound
		}
		return fmt.Errorf("event service: delete: %w", err)
	}

	return nil
}

func (e *eventServiceImpl) FindById(ctx context.Context, id uint) (*entity.EventResponse, error) {
	result, err := e.EventRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, exception.ErrorIdNotFound) {
			return nil, exception.ErrorIdNotFound
		}
		return nil, fmt.Errorf("event service: find by id: %w", err)
	}

	response := entity.ToEventResponse(result)

	return response, nil
}

func (e *eventServiceImpl) FindAll(ctx context.Context) ([]*entity.EventResponse, error) {
	result, err := e.EventRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("event service: find by all: %w", err)
	}

	var responses []*entity.EventResponse
	for _, v := range result {
		response := entity.EventResponse{
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			Capacity:    v.Capacity,
			Status:      v.Status,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}
