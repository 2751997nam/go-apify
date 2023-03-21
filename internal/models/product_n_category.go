package models

import "time"

type ProductNCategory struct {
	BaseModel
	ProductId  uint64    `json:"product_id"`
	CategoryId uint64    `json:"category_id"`
	IsParent   uint      `json:"is_parent"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ProductNCategory) TableName() string {
	return "sb_product_n_category"
}
