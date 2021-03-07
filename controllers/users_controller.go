package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ozzadar/platformer_mission_server/auth"
	"github.com/ozzadar/platformer_mission_server/models"
	"github.com/ozzadar/platformer_mission_server/responses"
	"github.com/ozzadar/platformer_mission_server/utils"
)

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	err = user.Validate("")

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	userCreated, err := user.NewUser(server.DBName, server.Session)

	if err != nil {
		formattedError := utils.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%s", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)
}

func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	users, err := user.FindAllUsers(server.DBName, server.Session)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := html.EscapeString(vars["id"])

	user := models.User{}

	userGotten, err := user.FindUserByID(server.DBName, server.Session, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, userGotten)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := html.EscapeString(vars["id"])

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

	tokenID, tokenEmail, err := auth.ExtractTokenIDAndEmail(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if (tokenID != id) && (tokenEmail != id) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedUser, err := user.UpdateAUser(server.DBName, server.Session, id)

	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, updatedUser)
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := html.EscapeString(vars["id"])

	user := models.User{}

	tokenID, tokenEmail, err := auth.ExtractTokenIDAndEmail(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if (tokenID != id) && (tokenEmail != id) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = user.DeleteAUser(server.DBName, server.Session, id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%s", id))
	responses.JSON(w, http.StatusNoContent, "")
}

func (server *Server) GetLevelsForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := html.EscapeString(vars["id"])

	user := models.User{}

	tokenID, tokenEmail, err := auth.ExtractTokenIDAndEmail(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if (tokenID != id) && (tokenEmail != id) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user.ID = tokenID

	levels, err := user.GetLevelIds(server.DBName, server.Session)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Getting levels for user successful",
		"levels":  levels,
	})
}
