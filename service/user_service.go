package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"ticket/domain/entity"
	"ticket/exception"
	"ticket/repository"
	"ticket/utils"
	"time"

	"github.com/go-playground/validator/v10"
)

type UserService interface {
	Create(ctx context.Context, req *entity.UserCreateRequest) (*entity.UserResponse, error)
	Update(ctx context.Context, userId uint, req *entity.UserUpdateRequest) (*entity.UserResponse, error)
	Delete(ctx context.Context, userId uint) error
	FindById(ctx context.Context, userId uint) (*entity.UserResponse, error)
	FindByEmail(ctx context.Context, email string) (*entity.UserResponse, error)
	FindAll(ctx context.Context) ([]*entity.UserResponse, error)
	Login(ctx context.Context, req *entity.UserLoginRequest) (*utils.TokenResponse, error)
	TokenRefresh(ctx context.Context, req *entity.UserRefreshTokenRequest) (*utils.TokenResponse, error)
}

type userServiceImpl struct {
	UserRepository repository.UserReposiitory
	Validate       *validator.Validate
}

func NewUserServiceImpl(userRepository repository.UserReposiitory, validate *validator.Validate) *userServiceImpl {
	return &userServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
	}
}

const (
	Admin    string = "admin"
	Customer string = "customer"
)

func (u *userServiceImpl) Create(ctx context.Context, req *entity.UserCreateRequest) (*entity.UserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if err := u.Validate.Struct(req); err != nil {
		return nil, exception.ErrorValidation
	}

	pass, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password %w", err)
	}

	user := entity.User{
		Name:     req.Name,
		Email:    email,
		Role:     Customer,
		Password: pass,
		Hp:       req.Hp,
		Address:  req.Address,
	}

	result, err := u.UserRepository.Create(ctx, &user)
	if err != nil {
		if errors.Is(err, exception.ErrorEmailExist) {
			return nil, exception.ErrorEmailExist
		}
		return nil, fmt.Errorf("user service: create: %w", err)
	}

	response := entity.ToUserResponse(result)

	return response, nil
}

func (u *userServiceImpl) Update(ctx context.Context, userId uint, req *entity.UserUpdateRequest) (*entity.UserResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		return nil, exception.ErrorValidation
	}

	user := entity.User{}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Password != nil {
		pass, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("hashing password %w", err)
		}

		user.Password = pass
	}

	if req.Hp != nil {
		user.Hp = *req.Hp
	}

	if req.Address != nil {
		user.Address = *req.Address
	}

	result, err := u.UserRepository.Update(ctx, userId, &user)
	if err != nil {
		return nil, fmt.Errorf("user service: update: %w", err)
	}

	response := entity.ToUserResponse(result)

	return response, nil
}

func (u *userServiceImpl) Delete(ctx context.Context, userId uint) error {
	if err := u.UserRepository.Delete(ctx, userId); err != nil {
		if errors.Is(err, exception.ErrorIdNotFound) {
			return exception.ErrorIdNotFound
		}
		return fmt.Errorf("user service: delete: %w", err)
	}

	return nil
}

func (u *userServiceImpl) FindById(ctx context.Context, userId uint) (*entity.UserResponse, error) {
	result, err := u.UserRepository.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, exception.ErrorIdNotFound) {
			return nil, exception.ErrorIdNotFound
		}
		return nil, fmt.Errorf("user service: find by id: %w", err)
	}

	response := entity.ToUserResponse(result)
	return response, nil
}

func (u *userServiceImpl) FindByEmail(ctx context.Context, email string) (*entity.UserResponse, error) {
	result, err := u.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, exception.ErrorEmailNotFound) {
			return nil, exception.ErrorEmailNotFound
		}
		return nil, fmt.Errorf("user service: find by email: %w", err)
	}

	response := entity.ToUserResponse(result)
	return response, nil
}

func (u *userServiceImpl) FindAll(ctx context.Context) ([]*entity.UserResponse, error) {
	result, err := u.UserRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("user service: find all: %w", err)
	}

	var responses []*entity.UserResponse
	for _, v := range result {
		response := entity.UserResponse{
			Name:      v.Name,
			Email:     v.Email,
			Hp:        v.Hp,
			Address:   v.Address,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}

func (u *userServiceImpl) Login(ctx context.Context, req *entity.UserLoginRequest) (*utils.TokenResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		return nil, exception.ErrorValidation
	}

	user, err := u.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, exception.ErrorEmailNotFound) {
			return nil, exception.ErrorFailedLogin
		}
		return nil, fmt.Errorf("user service: login find email: %w", err)
	}

	if !utils.CompareHashPassword(user.Password, req.Password) {
		return nil, exception.ErrorFailedLogin
	}

	tokenExp, _ := strconv.Atoi(os.Getenv("JWT_EXPIRED"))

	accessToken, err := utils.GenerateToken(user.ID, user.Name, user.Email, user.Role, time.Duration(tokenExp))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Name, user.Email, user.Role, time.Duration(tokenExp*2))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	createdToken := &utils.TokenResponse{
		Name:         user.Name,
		Token:        accessToken,
		TokenRefresh: refreshToken,
		TokenType:    "Bearer",
		ExipresIn:    tokenExp * 3600,
	}

	return createdToken, nil
}

func (u *userServiceImpl) TokenRefresh(ctx context.Context, req *entity.UserRefreshTokenRequest) (*utils.TokenResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		return nil, exception.ErrorValidation
	}

	tokenClaims, err := utils.ClaimTokenRefresh(req.TokenRefresh)
	if err != nil {
		return nil, exception.ErrorInvalidToken
	}

	user, err := u.UserRepository.FindById(ctx, tokenClaims.UserID)

	tokenExp, _ := strconv.Atoi(os.Getenv("JWT_EXPIRED"))

	accessToken, err := utils.GenerateToken(user.ID, user.Name, user.Email, user.Role, time.Duration(tokenExp))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	createdToken := &utils.TokenResponse{
		Name:         user.Name,
		Token:        accessToken,
		TokenRefresh: "",
		TokenType:    "Bearer",
		ExipresIn:    tokenExp * 3600,
	}

	return createdToken, nil
}
