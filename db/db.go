package db

import (
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(config *configs.Config) *gorm.DB {
	dialector := postgres.Open(config.Database.PostgresDSN)

	dbCon, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	dbCon.AutoMigrate(&models.User{})
	dbCon.AutoMigrate(&models.RefreshToken{})
	dbCon.AutoMigrate(&models.Log{})
	return dbCon
}
