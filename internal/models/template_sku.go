package models

type TemplateSku struct {
	BaseModel
	Sku        string             `json:"sku"`
	Price      float32            `json:"price"`
	HighPrice  float32            `json:"high_price"`
	ImageUrl   string             `json:"image_url"`
	IsDefault  int                `json:"is_default"`
	Status     string             `json:"status"`
	TemplateId uint64             `json:"template_id"`
	SkuValues  []TemplateSkuValue `gorm:"foreignKey:sku_id"`
	Gallery    []TemplateGallery  `gorm:"foreignKey:ProductId" json:"gallery,omitempty"`
}

func (TemplateSku) TableName() string {
	return "sb_product_template_sku"
}
