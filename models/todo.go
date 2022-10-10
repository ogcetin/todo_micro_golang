package models

type Todo struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	Detail       string `json:"detail"`
	CreationDate string `json:"creation_date"`
	LastUpdate   string `json:"last_update"`
	Status       string `json:"status"`
}
