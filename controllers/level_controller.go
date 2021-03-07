package controllers

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ozzadar/platformer_mission_server/auth"
	"github.com/ozzadar/platformer_mission_server/encrypt"
	"github.com/ozzadar/platformer_mission_server/models"
	"github.com/ozzadar/platformer_mission_server/responses"
)

func (server *Server) PostLevel(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	levelInformation := models.TileMap{}

	if err = json.Unmarshal(body, &levelInformation); err != nil {
		fmt.Println(levelInformation)
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	levelInformation.EncryptedMetadata = levelInformation.EncryptedMetadata[:len(levelInformation.EncryptedMetadata)-1]
	levelInformation.Metadata = *encrypt.DecryptMetaFile(levelInformation.EncryptedMetadata)

	if err = levelInformation.ValidateBeforeSave(); err != nil {
		fmt.Printf("problem4 %e\n", err)

		responses.ERROR(w, http.StatusPreconditionFailed, err)
		return
	}

	id, _, err := auth.ExtractTokenIDAndEmail(r)

	if err != nil {
		fmt.Printf("problem3 %e\n", err)

		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	levelInformation.Owner = id
	levelInformation.ID = levelInformation.Metadata.LocalID

	level, err := levelInformation.NewLevel(server.DBName, server.Session)

	if err != nil {
		fmt.Printf("problem2 %e\n", err)

		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Level create successfully",
		"level":   level,
	})
}

func (server *Server) GetLevel(w http.ResponseWriter, r *http.Request) {
	withMeta := true
	keys, ok := r.URL.Query()["withMeta"]

	log.Println(keys)

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'withMeta' is missing")
	} else {
		t, err := strconv.ParseBool(keys[0])

		if err != nil {
		} else {
			withMeta = t
		}
	}

	levelInformation := &models.TileMap{}

	vars := mux.Vars(r)
	id := html.EscapeString(vars["id"])

	log.Println(withMeta)
	levelInformation, err := levelInformation.GetLevel(id, withMeta, server.DBName, server.Session)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Success",
		"level":   levelInformation,
	})
}

func (server *Server) GetRandomLevels(w http.ResponseWriter, r *http.Request) {
	numberOfLevels := 5
	keys, ok := r.URL.Query()["count"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'withMeta' is missing")
	} else {
		requestNumber, err := strconv.Atoi(keys[0])

		if err != nil {
		} else {
			numberOfLevels = requestNumber
		}
	}

	levels, err := models.GetRandomSampleOfLevels(numberOfLevels, server.DBName, server.Session)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Success",
		"levels":  levels,
	})
}
