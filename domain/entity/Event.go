package entity

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID          uint   `gorm:"primaryKey;autoIncrement;notnull"`
	Name        string `gorm:"size:100"`
	Description string `gorm:"size:255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
