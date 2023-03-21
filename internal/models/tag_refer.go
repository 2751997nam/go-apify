package models

type TagRefer struct {
	BaseModel
	TagId     uint64 `json:"tag_id"`
	ReferId   uint64 `json:"refer_id"`
	ReferType string `json:"string"`
}

func (TagRefer) TableName() string {
	return "sb_tag_refer"
}
