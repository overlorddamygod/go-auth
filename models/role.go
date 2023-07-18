package models

type Role struct {
	ID   int    `gorm:"primary_key"`
	Name string `gorm:"column:name;unique" json:"name"`
	Time
}

func (Role) TableName() string {
	return "roles"
}
