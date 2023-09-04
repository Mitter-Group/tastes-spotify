package models

type TasteRequest struct {
	Name string `json:"name"`
}

type Taste struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DataResponse struct {
	UserId string  `json:"userId"`
	Data   []Taste `json:"data"`
	Source string  `json:"source" example:"SPOTIFY"`
}
