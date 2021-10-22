package inmemory

import (
	"context"
	"fmt"
	"lectures/hw6/internal/models"
	"lectures/hw6/internal/store"
	"sync"
)

type DB struct {
	data map[int]*models.Game

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		data: make(map[int]*models.Game),
		mu:   new(sync.RWMutex),
	}
}

func (db *DB) Create(ctx context.Context, game *models.Game) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[game.ID] = game
	return nil
}

func (db *DB) All(ctx context.Context) ([]*models.Game, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	games := make([]*models.Game, 0, len(db.data))
	for _, game := range db.data {
		games = append(games, game)
	}

	return games, nil
}

func (db *DB) ByID(ctx context.Context, id int) (*models.Game, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	game, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("No game with id %d", id)
	}

	return game, nil
}

func (db *DB) Update(ctx context.Context, game *models.Game) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[game.ID] = game
	return nil
}

func (db *DB) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
