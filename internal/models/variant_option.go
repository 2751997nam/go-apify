package models

import "time"

type VariantOption struct {
	ModelId
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Slug      string    `json:"slug"`
	VariantId uint64    `json:"variant_id"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `gorm:"column:created_at; type:timestamp; default: NOW(); not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; type:timestamp; default: NOW(); not null" json:"updated_at"`
}

func (VariantOption) TableName() string {
	return "sb_product_variant_option"
}

func (o *VariantOption) GetVariantId() uint64 {
	return o.VariantId
}
