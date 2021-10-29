package inmemory

import (
	"context"
	"fmt"
	"lectures/hw6/internal/models"
	"sync"
)

type GamesRepo struct {
	data map[int]*models.Game

	mu *sync.RWMutex
}

func (db *GamesRepo) Create(ctx context.Context, game *models.Game) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[game.ID] = game
	return nil
}

func (db *GamesRepo) All(ctx context.Context) ([]*models.Game, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	games := make([]*models.Game, 0, len(db.data))
	for _, game := range db.data {
		games = append(games, game)
	}

	return games, nil
}

func (db *GamesRepo) ByID(ctx context.Context, id int) (*models.Game, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	game, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("No game with id %d", id)
	}

	return game, nil
}

func (db *GamesRepo) Update(ctx context.Context, game *models.Game) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[game.ID] = game
	return nil
}

func (db *GamesRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
