package models

import "time"

type Variant struct {
	ModelId
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Variant) TableName() string {
	return "sb_product_variant"
}
