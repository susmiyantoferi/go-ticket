package handler

import (
	"errors"
	"net/http"
	"strconv"
	"ticket/domain/entity"
	"ticket/exception"
	"ticket/service"
	"ticket/web"

	"github.com/gin-gonic/gin"
)

type EventHandler interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindAll(ctx *gin.Context)
}

type eventHandlerImpl struct {
	EventService service.EventService
}

func NewEventHandlerImpl(eventService service.EventService) *eventHandlerImpl {
	return &eventHandlerImpl{
		EventService: eventService,
	}
}

func (e *eventHandlerImpl) Create(ctx *gin.Context) {
	req := entity.EventCreateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid request", nil)
		return
	}

	result, err := e.EventService.Create(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorValidation):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input", nil)
			return
		case errors.Is(err, exception.ErrorEventExist):
			web.ResponseJSON(ctx, http.StatusConflict, "error", "event already exist", nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusCreated, "success", "created", result)
}

func (e *eventHandlerImpl) Update(ctx *gin.Context) {
	req := entity.EventUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid request", nil)
		return
	}

	id := ctx.Param("id")
	eventId, err := strconv.Atoi(id)
	if err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input id type", nil)
		return
	}

	result, err := e.EventService.Update(ctx, uint(eventId), &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorValidation):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input", nil)
			return
		case errors.Is(err, exception.ErrorIdNotFound):
			web.ResponseJSON(ctx, http.StatusNotFound, "error", "id not found", nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "updated", result)
}

func (e *eventHandlerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	eventId, err := strconv.Atoi(id)
	if err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input id type", nil)
		return
	}

	if err := e.EventService.Delete(ctx, uint(eventId)); err != nil {
		switch {
		case errors.Is(err, exception.ErrorIdNotFound):
			web.ResponseJSON(ctx, http.StatusNotFound, "error", "id not found", nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "deleted", nil)
}

func (e *eventHandlerImpl) FindById(ctx *gin.Context) {
	id := ctx.Param("id")
	eventId, err := strconv.Atoi(id)
	if err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input id type", nil)
		return
	}

	result, err := e.EventService.FindById(ctx, uint(eventId))
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorIdNotFound):
			web.ResponseJSON(ctx, http.StatusNotFound, "error", "id not found", nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "success", result)
}

func (e *eventHandlerImpl) FindAll(ctx *gin.Context) {
	result, err := e.EventService.FindAll(ctx)
	if err != nil {
		web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
		return
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "success", result)
}
