package models

import "time"

type TemplateGalleryOverWrite struct {
	BaseModel
	ProductId    uint64    `json:"product_id"`
	ProductSkuId uint64    `json:"product_sku_id"`
	ImageUrl     string    `json:"image_url"`
	IsPrimary    int       `json:"is_primary"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (TemplateGalleryOverWrite) TableName() string {
	return "sb_product_template_gallery_overwrite"
}
