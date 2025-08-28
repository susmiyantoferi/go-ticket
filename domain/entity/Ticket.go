package entity

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	ID        uint    `gorm:"primaryKey;autoIncrement;notnull"`
	UserID    uint    `gorm:"notnull"`
	User      User    `gorm:"foreignKey:UserID;references:ID"`
	EventID   uint    `gorm:"notnull"`
	Event     Event   `gorm:"foreignKey:EventID;references:ID"`
	Qty       int     `gorm:"notnull"`
	UnitPrice float64 `gorm:"notnull"`
	Status    string  `gorm:"type:enum('confirm','canceled','waiting');default:'waiting';notnull"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
