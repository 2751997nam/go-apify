package models

import "time"

type TemplateGallery struct {
	BaseModel
	Type      string    `json:"type"`
	ProductId uint64    `json:"product_id"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TemplateGallery) TableName() string {
	return "sb_product_template_gallery"
}
