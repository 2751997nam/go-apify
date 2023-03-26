package models

import "time"

type BaseModel struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id,omitempty"`
}

type ModelTime struct {
	CreatedAt time.Time `gorm:"column:created_at; type:timestamp; default: NOW(); not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; type:timestamp; default: NOW(); not null" json:"updated_at"`
}
