package models

import "time"

type ProductGallery struct {
	ModelId
	Type      string `json:"type"`
	ProductId uint64 `json:"product_id"`
	ImageUrl  string `json:"image_url"`
	// Gallery   []string  `json:"gallery"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductGalleryData struct {
	ModelId
	Type      string    `json:"type"`
	ProductId uint64    `json:"product_id"`
	ImageUrl  string    `json:"image_url"`
	Gallery   []string  `json:"gallery"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductGallery) TableName() string {
	return "sb_product_gallery"
}
