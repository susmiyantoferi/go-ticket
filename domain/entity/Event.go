package entity

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID          uint    `gorm:"primaryKey;autoIncrement;notnull"`
	Name        string  `gorm:"size:100;unique;notnull"`
	Description string  `gorm:"size:255"`
	Price       float64 `gorm:"notnull"`
	Capacity    int     `gorm:"notnull"`
	Status      string  `gorm:"type:enum('aktif','berlangsung','selesai');default:'aktif';notnull"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type EventCreateRequest struct {
	Name        string  `validate:"required,min=1,max=100" json:"name"`
	Description string  `validate:"required,min=1,max=225" json:"description"`
	Price       float64 `validate:"required" json:"price"`
	Capacity    int     `validate:"required,gt=0" json:"capacity"`
}

type EventUpdateRequest struct {
	Name        *string  `validate:"omitempty,min=1,max=100" json:"name,omitempty"`
	Description *string  `validate:"omitempty,min=1,max=225" json:"description,omitempty"`
	Price       *float64 `validate:"omitempty" json:"price,omitempty"`
	Capacity    *int     `validate:"omitempty,gt=0" json:"capacity,omitempty"`
	Status      *string  `validate:"omitempty,oneof=aktif berlangsung selesai" json:"status,omitempty"`
}

type EventResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `son:"price"`
	Capacity    int       `json:"capacity"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToEventResponse(ev *Event) *EventResponse {
	return &EventResponse{
		Name:        ev.Name,
		Description: ev.Description,
		Price:       ev.Price,
		Capacity:    ev.Capacity,
		Status:      ev.Status,
		CreatedAt:   ev.CreatedAt,
		UpdatedAt:   ev.UpdatedAt,
	}
}
