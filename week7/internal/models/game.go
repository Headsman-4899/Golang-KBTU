package models

type Game struct {
	ID          int      `json:id`
	Name        string   `json:name`
	Description string   `json:description`
	Genre       []string `json:genre`

	Price     float64 `json:price`
	Developer string  `json:developer`
	Publisher string  `json:publisher`
}
