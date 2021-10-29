package inmemory

import (
	"context"
	"fmt"
	"lectures/hw6/internal/models"
	"sync"
)

type ProfilesRepo struct {
	data map[int]*models.Profile

	mu *sync.RWMutex
}

func (db *ProfilesRepo) Create(ctx context.Context, profile *models.Profile) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[profile.ID] = profile
	return nil
}

func (db *ProfilesRepo) All(ctx context.Context) ([]*models.Profile, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	profiles := make([]*models.Profile, 0, len(db.data))
	for _, profile := range db.data {
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (db *ProfilesRepo) ByID(ctx context.Context, id int) (*models.Profile, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	profile, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("No profile with id %d", id)
	}

	return profile, nil
}

func (db *ProfilesRepo) Update(ctx context.Context, profile *models.Profile) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[profile.ID] = profile
	return nil
}

func (db *ProfilesRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
