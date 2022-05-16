package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/db"
)

func Init() {
	r := NewRouter()
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	var err error
	if err = r.Run(":" + PORT); err != nil {
		log.Fatal(err)
	} else {
		log.Println("server started on port " + PORT)
	}
}

func InitForTest() *gin.Engine {
	configs.Load("../../.env")
	db.InitForTest()
	gin.SetMode(gin.ReleaseMode)
	return NewRouter()
}
