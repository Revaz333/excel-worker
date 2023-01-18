package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkerQueues struct {
	gorm.Model `json: "-"`
	ID         int            `json: "id", gorm:"primaryKey"`
	Driver     string         `json: "driver"`
	FilePath   string         `json: "file_path"`
	TableName  string         `json: "table_name"`
	Status     int            `json: "status"`
	DeletedAt  gorm.DeletedAt `json: "-"`
	CreatedAt  time.Time      `json: "created_at"`
	UpdatedAt  time.Time      `json: "updated_at"`
}
