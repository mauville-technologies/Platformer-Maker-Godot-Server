package models

import (
	"errors"
	"log"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type TileMap struct {
	ID                string          `rethinkdb:"id" json:"id,omitempty"`
	Owner             string          `rethinkdb:"user_id" json:"user_id,omitempty"`
	EncryptedMetadata []byte          `rethinkdb:"-" json:"encrypted_metadata,omitempty"`
	Metadata          Meta            `rethinkdb:"meta_data" json:"meta_data,omitempty"`
	Level             TileInformation `rethinkdb:"level" json:"tile_information,omitempty"`
}

type Meta struct {
	Completed      bool    `rethinkdb:"completed" json:"completed,omitempty"`
	LocalID        string  `rethinkdb:"id" json:"id,omitempty"`
	MapSizeX       int     `rethinkdb:"map_size_x" json:"map_size_x,omitempty"`
	MapSizeY       int     `rethinkdb:"map_size_y" json:"map_size_y,omitempty"`
	Name           string  `rethinkdb:"name" json:"name,omitempty"`
	StartPositionX int     `rethinkdb:"start_position_x" json:"start_position_x,omitempty"`
	StartPositionY int     `rethinkdb:"start_position_y" json:"start_position_y,omitempty"`
	OwnerTime      float64 `rethinkdb:"owner_time" json:"owner_time,omitempty"`
}

type TileInformation struct {
	TileList map[string]interface{} `rethinkdb:"tile_list" json:"tile_list,omitempty"`
	Tiles    map[string]interface{} `rethinkdb:"tiles" json:"tiles,omitempty"`
}

func (tm *TileMap) ValidateBeforeSave() error {
	if !tm.Metadata.Completed || tm.Metadata.OwnerTime == 0 {
		return errors.New("Level has not been properly completed before upload attempt")
	}

	return nil
}

func (tm *TileMap) NewLevel(dbname string, session *r.Session) (*TileMap, error) {

	resp, err := r.DB(dbname).Table("levels").Insert(tm, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(session)

	if err != nil || resp.Inserted < 1 {
		return nil, err
	}

	return tm, nil
}

func (tm *TileMap) GetLevel(id string, withMeta bool, dbname string, session *r.Session) (*TileMap, error) {
	request := r.DB(dbname).Table("levels").Filter(r.Row.Field("id").Eq(id))

	if withMeta {
		request = request.Pluck("id", "level", "user_id", "meta_data")
	} else {
		request = request.Pluck("id", "level", "user_id")
	}

	resp, err := request.Run(session)

	if err != nil {
		return nil, err
	}

	tilemap := &TileMap{}

	err = resp.One(tilemap)

	if err != nil {
		return nil, err
	}

	return tilemap, nil
}

func GetRandomSampleOfLevels(count int, dbname string, session *r.Session) ([]TileMap, error) {

	resp, err := r.DB(dbname).Table("levels").Sample(count).Run(session)

	if err != nil {
		return nil, err
	}

	levels := []TileMap{}

	err = resp.All(&levels)

	if err != nil {
		log.Println(levels, err)
		return nil, err
	}

	return levels, nil
}
