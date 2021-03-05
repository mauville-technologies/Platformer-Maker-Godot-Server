package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ozzadar/platformer_mission_server/auth"
	"github.com/ozzadar/platformer_mission_server/models"
	"github.com/ozzadar/platformer_mission_server/responses"
	"github.com/ozzadar/platformer_mission_server/utils"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	id, token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"token":   token,
		"user_id": id,
	})
}

func (server *Server) SignIn(email, password string) (string, string, error) {
	var err error

	user := models.User{}

	userToLogin, err := user.FindUserByEmail(server.DBName, server.Session, email)

	log.Print(userToLogin)
	if err != nil {
		return "", "", err
	}

	err = models.VerifyPassword(userToLogin.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", "", err
	}

	token, err := auth.CreateToken(userToLogin.ID, userToLogin.Email)

	return userToLogin.ID, token, err
}
