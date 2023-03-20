package models

import "time"

type ProductSku struct {
	ModelId
	Sku       string            `json:"sku"`
	Price     float32           `json:"price"`
	HighPrice float32           `json:"high_price"`
	ImageUrl  string            `json:"image_url"`
	ProductId uint64            `json:"product_id"`
	IsDefault int               `json:"is_default"`
	Status    string            `json:"status"`
	SkuValues []ProductSkuValue `gorm:"foreignKey:sku_id"`
	Gallery   []ProductGallery  `gorm:"foreignKey:ProductId" json:"gallery,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func (ProductSku) TableName() string {
	return "sb_product_sku"
}
