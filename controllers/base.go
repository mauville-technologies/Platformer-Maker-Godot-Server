package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ozzadar/platformer_mission_server/database"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type Server struct {
	DBName  string
	Session *r.Session
	Router  *mux.Router
}

func (server *Server) Initialize(DbContainerName, DbUser, DbPassword, DbHost, DbPort, DbName string, reseed bool) {

	// TODO: pass in DbUser, DbPassword into database.Init
	server.DBName = DbName
	if DbContainerName == "" {
		log.Println("No container name provided for DB, using HOST:PORT format")
		server.Session = database.Init(DbName, fmt.Sprintf("%s:%s", DbHost, DbPort), reseed)
	} else {
		log.Println("Container name provided for DB, using CONTAINER format")
		server.Session = database.Init(DbName, DbContainerName, reseed)
	}

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port " + addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
