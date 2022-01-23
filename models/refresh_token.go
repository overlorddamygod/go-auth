package models

import (
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	Token   string
	Revoked bool `gorm:"default:false"`
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
}
