package auth

import (
	"github.com/overlorddamygod/go-auth/mailer"
	"github.com/overlorddamygod/go-auth/models"
	"gorm.io/gorm"
)

type AuthController struct {
	db     *gorm.DB
	mailer *mailer.Mailer
	logger *models.Logger
}

func NewAuthController(db *gorm.DB, mailer *mailer.Mailer) AuthController {
	return AuthController{
		db:     db,
		mailer: mailer,
		logger: models.NewLogger(db),
	}
}
