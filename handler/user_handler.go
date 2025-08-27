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

type UserHandler interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindByEmail(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	Login(ctx *gin.Context)
	TokenRefresh(ctx *gin.Context)
}

type userHandlerImpl struct {
	UserService service.UserService
}

func NewUserHandlerImpl(userService service.UserService) *userHandlerImpl {
	return &userHandlerImpl{
		UserService: userService,
	}
}

func (u *userHandlerImpl) Create(ctx *gin.Context) {
	req := entity.UserCreateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid request", nil)
		return
	}

	result, err := u.UserService.Create(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorValidation):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input", nil)
			return
		case errors.Is(err, exception.ErrorEmailExist):
			web.ResponseJSON(ctx, http.StatusConflict, "error", "email already exist", nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusCreated, "success", "created", result)
}

func (u *userHandlerImpl) Update(ctx *gin.Context) {
	req := entity.UserUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid request", nil)
		return
	}

	userClaims, exist := ctx.Get("user")
	if !exist {
		web.ResponseJSON(ctx, http.StatusUnauthorized, "error", "user not found", nil)
		return
	}

	user := userClaims.(*utils.TokenClaim)

	userId := user.UserID

	result, err := u.UserService.Update(ctx, userId, &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorValidation):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input", nil)
			return
		case errors.Is(err, exception.ErrorIdNotFound):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "id not found", nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "updated", result)
}

func (u *userHandlerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("userId")
	userId, err := strconv.Atoi(id)
	if err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input type id", nil)
		return
	}

	if err := u.UserService.Delete(ctx, uint(userId)); err != nil {
		switch {
		case errors.Is(err, exception.ErrorIdNotFound):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "id not found", nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "deleted", nil)
}

func (u *userHandlerImpl) FindById(ctx *gin.Context) {
	userClaims, exist := ctx.Get("user")
	if !exist {
		web.ResponseJSON(ctx, http.StatusUnauthorized, "error", "user not found", nil)
		return
	}

	user := userClaims.(*utils.TokenClaim)

	userId := user.UserID

	result, err := u.UserService.FindById(ctx, userId)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorIdNotFound):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "id not found", nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "success", result)
}

func (u *userHandlerImpl) FindByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	result, err := u.UserService.FindByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorEmailNotFound):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "email not found", nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "success", result)
}

func (u *userHandlerImpl) FindAll(ctx *gin.Context) {
	result, err := u.UserService.FindAll(ctx)
	if err != nil {
		switch {
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "success", result)
}

func (u *userHandlerImpl) Login(ctx *gin.Context) {
	req := entity.UserLoginRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid request", nil)
		return
	}

	result, err := u.UserService.Login(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorValidation):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input", nil)
			return
		case errors.Is(err, exception.ErrorFailedLogin):
			web.ResponseJSON(ctx, http.StatusUnauthorized, "unauthorizhed", err.Error(), nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, "error", "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "success", result)
}

func (u *userHandlerImpl) TokenRefresh(ctx *gin.Context) {
	req := entity.UserRefreshTokenRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid request", nil)
		return
	}

	result, err := u.UserService.TokenRefresh(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrorValidation):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid input", nil)
			return
		case errors.Is(err, exception.ErrorInvalidToken):
			web.ResponseJSON(ctx, http.StatusBadRequest, "error", "invalid token input", nil)
			return
		default:
			web.ResponseJSON(ctx, http.StatusInternalServerError, err.Error(), "something wrong", nil)
			return
		}
	}

	web.ResponseJSON(ctx, http.StatusOK, "success", "success", result)
}
