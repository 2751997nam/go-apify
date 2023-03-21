package models

import "fmt"

type ProductNTemplate struct {
	BaseModel
	ProductId  uint64 `json:"product_id"`
	TemplateId uint64 `json:"template_id"`
}

func (ProductNTemplate) TableName() string {
	return "sb_product_n_template"
}

func (p *ProductNTemplate) Exists(productId uint64) bool {
	db := GetDB()
	var Found bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM `%s` WHERE product_id = ?) AS found", p.TableName())
	db.Raw(query, productId).Scan(&Found)

	return Found
}
