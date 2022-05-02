package models

import (
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	Token     string
	Revoked   bool `gorm:"default:false"`
	UserAgent string
	IP        string
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
}
