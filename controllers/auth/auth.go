package auth

import (
	"github.com/overlorddamygod/go-auth/mailer"
	"gorm.io/gorm"
)

type AuthController struct {
	db     *gorm.DB
	mailer *mailer.Mailer
}

func NewAuthController(db *gorm.DB, mailer *mailer.Mailer) AuthController {
	return AuthController{
		db:     db,
		mailer: mailer,
	}
}
