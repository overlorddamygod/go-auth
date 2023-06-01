package admin

import (
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/mailer"
	"github.com/overlorddamygod/go-auth/models"
	"gorm.io/gorm"
)

type AdminController struct {
	config *configs.Config
	db     *gorm.DB
	mailer *mailer.Mailer
	logger *models.Logger
}

func NewAdminController(config *configs.Config, db *gorm.DB, mailer *mailer.Mailer, logger *models.Logger) *AdminController {
	return &AdminController{
		config: config,
		db:     db,
		mailer: mailer,
		logger: logger,
	}
}
