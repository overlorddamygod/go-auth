package models

import (
	"github.com/google/uuid"
)

type RefreshToken struct {
	Basic
	Token     string
	Revoked   bool `gorm:"default:false"`
	UserAgent string
	IP        string
	UserID    uuid.UUID
	User      User `gorm:"foreignKey:UserID"`
}
