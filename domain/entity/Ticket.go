package entity

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	ID          uint    `gorm:"primaryKey;autoIncrement;notnull"`
	UserID      uint    `gorm:"notnull"`
	User        User    `gorm:"foreignKey:UserID;references:ID"`
	EventID     uint    `gorm:"notnull"`
	Event       Event   `gorm:"foreignKey:EventID;references:ID"`
	Qty         int     `gorm:"notnull"`
	UnitPrice   float64 `gorm:"notnull"`
	TotalAmount float64 `gorm:"notnull"`
	Status      string  `gorm:"type:enum('confirm','canceled','waiting');default:'waiting';notnull"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type TicketCreateRequest struct {
	UserID  uint `validate:"required" json:"user_id"`
	EventID uint `validate:"required" json:"event_id"`
	Qty     int  `validate:"required,gt=0" json:"qty"`
}

type TicketUpdateRequest struct {
	Status string `validate:"required,oneof=confirm canceled waiting" json:"status"`
}

type UserInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Hp      string `json:"hp"`
	Address string `json:"address"`
}

type EventInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TicketResponse struct {
	UserID      uint      `json:"user_id"`
	User        UserInfo  `json:"user"`
	EventID     uint      `json:"event_id"`
	Event       EventInfo `json:"event"`
	Qty         int       `json:"qty"`
	UnitPrice   float64   `json:"unit_price"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status_ticket"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToTicketResponse(ticks *Ticket) *TicketResponse {
	return &TicketResponse{
		UserID: ticks.UserID,
		User: UserInfo{
			Email:   ticks.User.Email,
			Name:    ticks.User.Name,
			Hp:      ticks.User.Hp,
			Address: ticks.User.Address,
		},
		EventID: ticks.EventID,
		Event: EventInfo{
			Name:        ticks.Event.Name,
			Description: ticks.Event.Description,
		},
		Qty:         ticks.Qty,
		UnitPrice:   ticks.UnitPrice,
		TotalAmount: ticks.TotalAmount,
		Status:      ticks.Status,
		CreatedAt:   ticks.CreatedAt,
		UpdatedAt:   ticks.UpdatedAt,
	}
}
