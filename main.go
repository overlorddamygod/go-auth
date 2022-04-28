package main

import (
	"github.com/overlorddamygod/go-auth/db"
	"github.com/overlorddamygod/go-auth/server"
)

func main() {
	db.Init()
	server.Init()
}
