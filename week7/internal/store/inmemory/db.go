package inmemory

import (
	"context"
	"errors"
	"fmt"
	"lectures/hw6/internal/models"
	"lectures/hw6/internal/store"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	gamesRepo *mongo.Collection
	//profilesRepo *mongo.Collection

	mu *sync.RWMutex
}

func (db *DB) Create(ctx context.Context, game *models.Game) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.gamesRepo.InsertOne(ctx, game)
	return err
}

func (db *DB) All(ctx context.Context, filter *models.GamesFilter) ([]*models.Game, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	//filter := bson.D{{}}

	return db.filterTasks(ctx, filter)
}

func (db *DB) ByID(ctx context.Context, id string) (*models.Game, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("conversion of id from string to ObjectID: %s", id)
	}
	filter := bson.D{primitive.E{Key: "_id", Value: idObj}}
	u := &models.Game{}
	ok := db.gamesRepo.FindOne(ctx, filter).Decode(u)
	if ok != nil {
		return nil, fmt.Errorf("no user with id %s", id)
	}
	return u, nil
}

func (db *DB) Update(ctx context.Context, game *models.Game) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	filter := bson.D{primitive.E{Key: "_id", Value: game.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "name", Value: game.Name},
		primitive.E{Key: "description", Value: game.Description},
		primitive.E{Key: "genre", Value: game.Genre},
		primitive.E{Key: "price", Value: game.Price},
		primitive.E{Key: "developer", Value: game.Developer},
		primitive.E{Key: "publisher", Value: game.Publisher},
		primitive.E{Key: "reviews", Value: game.Reviews},
	}}}

	u := &models.Game{}

	return db.gamesRepo.FindOneAndUpdate(ctx, filter, update).Decode(u)
}

func (db *DB) Delete(ctx context.Context, id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("conversion of id from string to ObjectID: %s", id)
	}
	filter := bson.D{primitive.E{Key: "_id", Value: idObj}}

	res, err := db.gamesRepo.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no tasks were deleted")
	}

	return nil
}

func Init() store.GamesRepository {
	var ctx = context.TODO()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	return &DB{
		gamesRepo: client.Database("GamesDB").Collection("games"),
		//profilesRepo: client.Database("ProfilesDB").Collection("profiles"),
		mu: new(sync.RWMutex),
	}
}

func (db *DB) filterTasks(ctx context.Context, filter interface{}) ([]*models.Game, error) {
	var games []*models.Game

	cur, err := db.gamesRepo.Find(ctx, filter)
	if err != nil {
		return games, err
	}

	for cur.Next(ctx) {
		var t models.Game
		err := cur.Decode(&t)
		if err != nil {
			return games, err
		}

		games = append(games, &t)
	}

	if err := cur.Err(); err != nil {
		return games, err
	}

	cur.Close(ctx)

	if len(games) == 0 {
		return games, mongo.ErrNoDocuments
	}

	return games, nil
}
