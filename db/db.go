package db

import (
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/models"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	dialector, err := configs.GetConfig().Database.GetDialector()

	if err != nil {
		panic(err)
	}

	dbCon, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	db = dbCon
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.RefreshToken{})
}

func GetDB() *gorm.DB {
	return db
}
