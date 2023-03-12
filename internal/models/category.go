package models

import (
	"database/sql"
	"time"
)

type Category struct {
	ID            int           `gorm:"PRIMARY_KEY;AUTO_INCREMENT" nestedset:"id" json:"id"`
	Name          string        `nestedset:"name" json:"name"`
	Slug          string        `nestedset:"slug" json:"slug"`
	Description   string        `nestedset:"description" json:"description"`
	ImageUrl      string        `nestedset:"image_url" json:"image_url"`
	BigImageUrl   string        `nestedset:"big_image_url" json:"big_image_url"`
	Type          string        `nestedset:"type" json:"type"`
	Breadcrumb    string        `nestedset:"breadcrumb" json:"breadcrumb"`
	Lft           int           `nestedset:"lft" json:"_lft" gorm:"column:_lft"`
	Rgt           int           `nestedset:"rgt" json:"_rgt" gorm:"column:_rgt"`
	Depth         int           `nestedset:"depth" json:"depth" gorm:"column:depth"`
	ParentId      sql.NullInt64 `nestedset:"parent_id" json:"parent_id" gorm:"column:parent_id"`
	ChildrenCount int           `nestedset:"children_count" json:"children_count" gorm:"column:children_count"`
	IsHidden      int           `nestedset:"is_hidden" json:"is_hidden"`
	CreatedAt     time.Time     `nestedset:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `nestedset:"updated_at" json:"updated_at"`
}

func (Category) TableName() string {
	return "sb_category"
}
