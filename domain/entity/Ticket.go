package entity

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;notnull"`
	EventID   uint   `gorm:"notnull"`
	Event     Event  `gorm:"foreignKey:EventID;references:ID"`
	Code      string `gorm:"size:100;notnull"`
	Status    string `gorm:"type:enum('confirm','canceled','waiting');default:'waiting';notnull"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
