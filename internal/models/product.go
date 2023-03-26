package models

import "time"

type Product struct {
	BaseModel
	Sku                string           `json:"sku,omitempty"`
	Name               string           `json:"name,omitempty"`
	Slug               string           `json:"slug,omitempty"`
	Price              float32          `json:"price,omitempty"`
	HighPrice          float32          `json:"high_price,omitempty"`
	AddShippingFee     float32          `json:"add_shipping_fee,omitempty"`
	ImageUrl           string           `json:"image_url,omitempty"`
	Status             string           `json:"status,omitempty"`
	Description        string           `json:"description,omitempty"`
	Content            string           `json:"content,omitempty"`
	BrandId            uint64           `json:"brand_id,omitempty"`
	Note               string           `json:"note,omitempty"`
	ApproveAdvertising int              `json:"approve_advertising,omitempty"`
	IsTrademark        int              `json:"is_trademark,omitempty"`
	ActorId            uint64           `json:"actor_id,omitempty"`
	UpdaterId          uint64           `json:"updater_id,omitempty"`
	CreatedAt          time.Time        `json:"created_at,omitempty"`
	UpdatedAt          time.Time        `json:"updated_at,omitempty"`
	Gallery            []ProductGallery `gorm:"foreignKey:ProductId" json:"gallery,omitempty"`
	Creater            User             `gorm:"foreignKey:ActorId" json:"creator,omitempty"`
	Modifier           User             `gorm:"foreignKey:UpdaterId" json:"modifier,omitempty"`
	Tags               []Tag            `gorm:"many2many:sb_tag_refer;joinForeignKey:ReferId;joinReferences:TagId" json:"tags,omitempty"`
	Categories         []Category       `gorm:"many2many:sb_product_n_category;joinForeignKey:ProductId;joinReferences:CategoryId" json:"categories,omitempty"`
}

func (Product) TableName() string {
	return "sb_product"
}
