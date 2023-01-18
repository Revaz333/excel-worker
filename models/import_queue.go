package models

import (
	"time"

	"gorm.io/gorm"
)

type ImportQueue struct {
	gorm.Model `json: "-"`
	ID         int            `json: "id", gorm:"primaryKey"`
	Query      string         `json: "query"`
	DeletedAt  gorm.DeletedAt `json: "-"`
	CreatedAt  time.Time      `json: "created_at"`
	UpdatedAt  time.Time      `json: "updated_at"`
}
