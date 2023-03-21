package models

import "time"

type TemplateDesignCode struct {
	BaseModel
	DesignCode string    `json:"design_code"`
	ProductId  uint64    `json:"product_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (TemplateDesignCode) TableName() string {
	return "sb_product_template_design_code"
}
