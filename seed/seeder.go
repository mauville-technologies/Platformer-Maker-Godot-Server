package seed

import (
	"log"

	"github.com/ozzadar/platformer_mission_server/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var users = []models.User{
	models.User{
		Nickname: "Demo",
		Email:    "demo@demo.com",
		Password: "password",
	},
}

func Load(dbName string, session *r.Session) {
	for _, user := range users {
		if _, err := user.NewUser(dbName, session); err != nil {
			log.Fatalf("cannot seed users table: %v: ", err)
		}
	}
}
