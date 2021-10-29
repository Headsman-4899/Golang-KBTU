package models

type Profile struct {
	ID       int    `json:id`
	Nickname string `json:nickname`
	Country  string `json:country`
	Level    int    `json:level`

	Friends      []Profile `json:friends`
	GamesLibrary []Game    `json: gamesLibrary`
	WishList     []Game    `json: wishlist`
}
