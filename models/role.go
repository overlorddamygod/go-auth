package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        int    `primary_key`
	Name      string `gorm:"column:name" json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Role) TableName() string {
	return "roles"
}
