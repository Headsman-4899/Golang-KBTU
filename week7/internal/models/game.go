package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Game struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Genre       []string           `bson:"genre"`

	Price     float64  `bson:"price"`
	Developer string   `bson:"developer"`
	Publisher string   `bson:"publisher"`
	Reviews   []string `bson:"reviews"`
}

type GamesFilter struct {
	Query *string `json:"query"`
}
