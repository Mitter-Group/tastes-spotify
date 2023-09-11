package models

import (
	"fmt"
)

type DataType string

const (
	Tracks  DataType = "tracks"
	Artists DataType = "artists"
	Genres  DataType = "genres"
)

type DataRequest struct {
	DataType DataType `json:"data_type" example:"tracks"`
	UserId   string   `json:"user_id" example:"123456789"`
}

type Data struct {
	UserId string        `json:"user_id"`
	Data   []DataDetails `json:"data"`
	Source string        `json:"source" example:"SPOTIFY"`
}

type DataDetails struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (dr *DataRequest) Validate() error {
	switch dr.DataType {
	case Tracks, Artists, Genres:
		return nil
	default:
		return fmt.Errorf("Invalid DataType value: %s", dr.DataType)
	}
}
