package models

type User struct {
	BaseModel
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func (User) TableName() string {
	return "sb_users"
}
