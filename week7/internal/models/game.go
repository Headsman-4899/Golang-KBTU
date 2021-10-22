package models

type Game struct {
	ID    int    `json:id`
	Name  string `json:name`
	Genre string `json:genre`
}
