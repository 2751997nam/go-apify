package models

import "time"

type ProductSkuValue struct {
	ModelId
	ProductId       uint64    `json:"product_id"`
	SkuId           uint64    `json:"sku_id"`
	VariantId       uint64    `json:"variant_id"`
	VariantOptionId uint64    `json:"variant_option_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (ProductSkuValue) TableName() string {
	return "sb_product_sku_value"
}
