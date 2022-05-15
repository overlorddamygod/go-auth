package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Basic struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
