package types

type ProductSkuView struct {
	ID uint64 `json:"id"`
	// Sku       string  `json:"sku"`
	Price     float32 `json:"price"`
	HighPrice float32 `json:"high_price"`
	// ImageUrl  string       `json:"image_url"`
	ProductId uint64 `json:"product_id"`
	IsDefault int    `json:"is_default"`
	Status    string `json:"status"`
	// Variants  []OptionView `json:"variants"`
	Variants []uint64 `json:"variants"`
	Gallery  []string `json:"gallery"`
}
