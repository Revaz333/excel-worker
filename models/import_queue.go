package models

type ImportQueue struct {
	ID    int    `json: "id", gorm:"primaryKey"`
	Query string `json: "query"`
}
