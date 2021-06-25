package main

import (
	"log"

	"github.com/ash3798/Social-Network/config"
	"github.com/ash3798/Social-Network/database"
	"github.com/ash3798/Social-Network/server"
)

func main() {
	log.Println("social network")

	config.InitEnv()

	err := database.InitDatabase()
	if err != nil {
		return
	}

	err = database.Action.PrepareDatabase()
	if err != nil {
		return
	}

	config.InitReactions()

	server.StartServer()

	database.Action.CloseDBConnection()
}
