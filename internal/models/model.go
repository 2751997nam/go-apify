package models

import "time"

type ModelId struct {
	ID uint `gorm:"primary_key" json:"id"`
}

type ModelTime struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
