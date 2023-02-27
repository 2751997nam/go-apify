package models

type ProductSku struct {
	ModelId
	Sku       string `json:"sku"`
	Price     string `json:"price"`
	HighPrice string `json:"high_price"`
	ImageUrl  string `json:"image_url"`
	ProductId string `json:"product_id"`
	IsDefault string `json:"is_default"`
	ModelTime
}

func (ProductSku) TableName() string {
	return "sb_product_sku"
}
