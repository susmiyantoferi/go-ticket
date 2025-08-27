package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:100;notnull"`
	Role      string `gorm:"type:enum('admin','customer');default:'customer';notnull"`
	Email     string `gorm:"size:100;unique;notnull"`
	Password  string `gorm:"size:255;notnull"`
	Hp        string `gorm:"size:20,notnull"`
	Address   string `gorm:"notnull"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserCreateRequest struct {
	Name     string `validate:"required,min=1,max=100" json:"name"`
	Email    string `validate:"required,email,min=1,max=100" json:"email"`
	Password string `validate:"required,min=1,max=255" json:"password"`
	Hp       string `validate:"required,min=1,max=100" json:"hp"`
	Address  string `validate:"required" json:"address"`
}

type UserUpdateRequest struct {
	Name     *string `validate:"omitempty,min=1,max=100" json:"name,omitempty"`
	Password *string `validate:"omitempty,min=1,max=255" json:"password,omitempty"`
	Hp       *string `validate:"omitempty,min=1,max=20,numeric" json:"hp,omitempty"`
	Address  *string `validate:"omitempty" json:"address,omitempty"`
}

type UserResponse struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Hp        string    `json:"hp"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLoginRequest struct {
	Email string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
}

type UserRefreshTokenRequest struct {
	TokenRefresh string `validate:"required" json:"token_refresh"`
}

func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		Hp:        user.Hp,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
