package server

import (
	"log"
	"os"
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
