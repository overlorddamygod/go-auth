package main

import (
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/db"
	"github.com/overlorddamygod/go-auth/server"
)

func main() {
	configs.LoadConfig()
	db.Init()
	server.Init()
}
