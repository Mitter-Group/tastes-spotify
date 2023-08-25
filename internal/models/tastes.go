package models

type TasteRequest struct {
	Name string `json:"name"`
}

type Taste struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
