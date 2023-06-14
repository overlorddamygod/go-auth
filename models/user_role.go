package models

import (
	"github.com/google/uuid"
)

type UserRole struct {
	UserID uuid.UUID `gorm:"type:uuid;column:user_id" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID"`
	Type   int       `gorm:"column:type" json:"type"`
	Role   Role      `gorm:"foreignKey:Type"`
	Basic
}

func (UserRole) TableName() string {
	return "user_roles"
}
