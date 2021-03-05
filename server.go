package main

import (
	"os"

	"github.com/ozzadar/platformer_mission_server/controllers"
	"github.com/ozzadar/platformer_mission_server/seed"
)

var Server = controllers.Server{}

func Run(reseed bool) {
	Server.Initialize(os.Getenv("DB_CONTAINER_NAME"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), reseed)

	if reseed {
		seed.Load(Server.DBName, Server.Session)
	}

	Server.Run(":" + os.Getenv("SERVER_PORT"))
}
