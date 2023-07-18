package models

import (
	"github.com/google/uuid"
)

type UserRole struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID"`
	Type   int       `gorm:"primaryKey;column:type" json:"type"`
	Role   Role      `gorm:"foreignKey:Type"`
	Time
}

func (UserRole) TableName() string {
	return "user_roles"
}
