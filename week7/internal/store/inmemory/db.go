package inmemory

import (
	"lectures/hw6/internal/models"
	"lectures/hw6/internal/store"
	"sync"
)

type DB struct {
	gamesRepo    store.GamesRepository
	profilesRepo store.ProfilesRepository

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		mu: new(sync.RWMutex),
	}
}

func (db *DB) Games() store.GamesRepository {
	if db.gamesRepo == nil {
		db.gamesRepo = &GamesRepo{
			data: make(map[int]*models.Game),
			mu:   new(sync.RWMutex),
		}
	}

	return db.gamesRepo
}

func (db *DB) Profiles() store.ProfilesRepository {
	if db.profilesRepo == nil {
		db.profilesRepo = &ProfilesRepo{
			data: make(map[int]*models.Profile),
			mu:   new(sync.RWMutex),
		}
	}

	return db.profilesRepo
}
