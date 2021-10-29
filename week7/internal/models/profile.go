package models

type Profile struct {
	ID            int    `json:id`
	Nickname      string `json:nickname`
	Country       string `json:country`
	Level         int    `json:level`
	GamesQuantity int    `json: gamesQuantity`
	WishList      []Game `json: wishlist`
}
