package models

type Template struct {
	BaseModel
	Name          string `json:"name"`
	ProductIdFake uint64 `json:"product_id_fake"`
}

func (Template) TableName() string {
	return "sb_product_template"
}
