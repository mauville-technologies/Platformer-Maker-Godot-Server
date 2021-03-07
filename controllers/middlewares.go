package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/ozzadar/platformer_mission_server/auth"
	"github.com/ozzadar/platformer_mission_server/models"
	"github.com/ozzadar/platformer_mission_server/responses"
)

func (s *Server) SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Access-Control-Allow-Origin", "*")
   w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
   w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
if (*r).Method == "OPTIONS" {
      return
   }
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func (s *Server) SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)

		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthenticated"))
			return
		}

		id, _, err := auth.ExtractTokenIDAndEmail(r)

		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthenticated"))
			return
		}

		user := &models.User{}

		if user, err = user.FindUserByID(s.DBName, s.Session, id); err != nil {
			log.Println("User in token doesn't exist.")
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthenticated"))
			return
		}

		next(w, r)
	}
}
