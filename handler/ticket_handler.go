package handler

import (
	"errors"
	"net/http"
	"strconv"
	"ticket/domain/entity"
	"ticket/exception"
	"ticket/service"
	"ticket/utils"
	"ticket/web"

	"github.com/gin-gonic/gin"
)

type TicketHandler interface {
	Create(cxt *gin.Context)
	Update(cxt *gin.Context)
	Delete(cxt *gin.Context)
	FindById(cxt *gin.Context)
	FindByUserId(cxt *gin.Context)
	FindAll(cxt *gin.Context)
}

type ticketHandlerImpl struct {
	TicketService service.TicketService
}

func NewTicketHandlerImpl(ticketService service.TicketService) *ticketHandlerImpl {
	return &ticketHandlerImpl{
		TicketService: ticketService,
	}
}

func (t *ticketHandlerImpl) Create(ctx *gin.Context) {
	req := entity.TicketCreateRequest{}

	userlaims, exist := ctx.Get("user")
	if !exist {
		web.ResponseJSON(ctx, http.StatusUnauthorized, "error", "user not found", nil)
		return
	}

	user := userlaims.(*utils.TokenClaim)
	req.UserID = user.UserID

	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid request", nil)
		return
	}

	result, err := t.TicketService.Create(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorValidation):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input", nil)
			return
		case errors.Is(err, exception.ErrorEventNotFound):
			web.ResponseJSON(ctx, http.StatusNotFound, "error", err.Error(), nil)
			return
		case errors.Is(err, exception.ErrorStockNotEnough):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", err.Error(), nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
			return
		}
	}
	web.ResponseJSON(ctx, http.StatusOK, "success", "created", result)
}

func (t *ticketHandlerImpl) Update(ctx *gin.Context) {
	req := entity.TicketUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid request", nil)
		return
	}

	id := ctx.Param("id")
	ticksId, err := strconv.Atoi(id)
	if err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input id type", nil)
		return
	}

	result, err := t.TicketService.Update(ctx, uint(ticksId), &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorValidation):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input", nil)
			return
		case errors.Is(err, exception.ErrorIdNotFound):
			web.ResponseJSON(ctx, http.StatusNotFound, "error", err.Error(), nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "updated", result)
}

func (t *ticketHandlerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	ticksId, err := strconv.Atoi(id)
	if err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input id type", nil)
		return
	}

	if err := t.TicketService.Delete(ctx, uint(ticksId)); err != nil {
		switch {
		case errors.Is(err, exception.ErrorIdNotFound):
			web.ResponseJSON(ctx, http.StatusNotFound, "error", err.Error(), nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "deleted", nil)
}

func (t *ticketHandlerImpl) FindById(ctx *gin.Context) {
	id := ctx.Param("id")
	ticksId, err := strconv.Atoi(id)
	if err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input id type", nil)
		return
	}

	result, err := t.TicketService.FindById(ctx, uint(ticksId))
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorIdNotFound):
			web.ResponseJSON(ctx, http.StatusNotFound, "error", err.Error(), nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "success", result)
}

func (t *ticketHandlerImpl) FindByUserId(ctx *gin.Context) {
	userlaims, exist := ctx.Get("user")
	if !exist {
		web.ResponseJSON(ctx, http.StatusUnauthorized, "error", "user not found", nil)
		return
	}

	user := userlaims.(*utils.TokenClaim)
	userId := user.UserID

	result, err := t.TicketService.FindByUserId(ctx, userId)
	if err != nil {
		switch {
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "success", result)
}

func (t *ticketHandlerImpl) FindAll(ctx *gin.Context) {
	result, err := t.TicketService.FindAll(ctx)
	if err != nil {
		switch {
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "success", result)
}
