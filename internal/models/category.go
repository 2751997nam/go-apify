package models

import (
	"database/sql"
	"time"
)

type Category struct {
	ID            int           `gorm:"PRIMARY_KEY;AUTO_INCREMENT" nestedset:"id" json:"id"`
	Name          string        `nestedset:"name" json:"name,omitempty"`
	Slug          string        `nestedset:"slug" json:"slug,omitempty"`
	Description   string        `nestedset:"description" json:"description,omitempty"`
	ImageUrl      string        `nestedset:"image_url" json:"image_url,omitempty"`
	BigImageUrl   string        `nestedset:"big_image_url" json:"big_image_url,omitempty"`
	Type          string        `nestedset:"type" json:"type,omitempty"`
	Breadcrumb    string        `nestedset:"breadcrumb" json:"breadcrumb,omitempty" grom:"column:breadcrumb"`
	Lft           int           `nestedset:"lft" json:"_lft,omitempty" gorm:"column:_lft"`
	Rgt           int           `nestedset:"rgt" json:"_rgt,omitempty" gorm:"column:_rgt"`
	Depth         int           `nestedset:"depth" json:"depth" gorm:"column:depth,omitempty"`
	ParentId      sql.NullInt64 `nestedset:"parent_id" json:"parent_id,omitempty" gorm:"column:parent_id"`
	ChildrenCount int           `nestedset:"children_count" json:"children_count,omitempty" gorm:"column:children_count"`
	IsHidden      int           `nestedset:"is_hidden" json:"is_hidden,omitempty"`
	CreatedAt     *time.Time    `nestedset:"created_at" json:"created_at,omitempty"`
	UpdatedAt     *time.Time    `nestedset:"updated_at" json:"updated_at,omitempty"`
}

func (Category) TableName() string {
	return "sb_category"
}
