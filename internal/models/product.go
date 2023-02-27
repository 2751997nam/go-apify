package models

type Product struct {
	ModelId
	Sku                string  `json:"sku"`
	Name               string  `json:"name"`
	Slug               string  `json:"slug"`
	Price              float32 `json:"price"`
	HighPrice          float32 `json:"high_price"`
	AddShippingFee     float32 `json:"add_shipping_fee"`
	ImageUrl           string  `json:"image_url"`
	Status             string  `json:"status"`
	Description        string  `json:"description"`
	Content            string  `json:"content"`
	BrandId            uint    `json:"brand_id"`
	Note               string  `json:"note"`
	ApproveAdvertising int     `json:"approve_advertising"`
	IsTrademark        int     `json:"is_trademark"`
	ModelTime
}

func (Product) TableName() string {
	return "sb_product"
}
