package db

import (
	"github.com/overlorddamygod/go-auth/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	dbConn, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = dbConn
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.RefreshToken{})
}

func GetDB() *gorm.DB {
	return db
}
