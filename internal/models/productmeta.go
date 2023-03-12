package models

import "time"

type ProductMeta struct {
	ModelId
	ProductId uint64    `json:"product_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductMeta) TableName() string {
	return "sb_product_meta"
}
