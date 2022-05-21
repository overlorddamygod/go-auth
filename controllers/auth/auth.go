package auth

import (
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/mailer"
	"github.com/overlorddamygod/go-auth/models"
	"gorm.io/gorm"
)

type AuthController struct {
	config *configs.Config
	db     *gorm.DB
	mailer *mailer.Mailer
	logger *models.Logger
}

func NewAuthController(config *configs.Config, db *gorm.DB, mailer *mailer.Mailer, logger *models.Logger) *AuthController {
	return &AuthController{
		config: config,
		db:     db,
		mailer: mailer,
		logger: logger,
	}
}
