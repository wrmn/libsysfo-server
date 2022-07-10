package main

import (
	"libsysfo-server/database"
	"libsysfo-server/server"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	} else {
		go database.Checker()
		server.Serve(port)
	}
}
