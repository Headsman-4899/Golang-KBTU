package store

import (
	"context"
	"lectures/hw6/internal/models"
)

type Store interface {
	Games() GamesRepository
	Profiles() ProfilesRepository
}

type GamesRepository interface {
	Create(ctx context.Context, game *models.Game) error
	All(ctx context.Context) ([]*models.Game, error)
	ByID(ctx context.Context, id int) (*models.Game, error)
	Update(ctx context.Context, game *models.Game) error
	Delete(ctx context.Context, id int) error
}

type ProfilesRepository interface {
	Create(ctx context.Context, profile *models.Profile) error
	All(ctx context.Context) ([]*models.Profile, error)
	ByID(ctx context.Context, id int) (*models.Profile, error)
	Update(ctx context.Context, profile *models.Profile) error
	Delete(ctx context.Context, id int) error
}
