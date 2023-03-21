package models

type TemplateSkuValue struct {
	BaseModel
	SkuId           uint64 `json:"sku_id"`
	VariantId       uint64 `json:"variant_id"`
	VariantOptionId uint64 `json:"variant_option_id"`
}

func (TemplateSkuValue) TableName() string {
	return "sb_product_template_sku_value"
}
