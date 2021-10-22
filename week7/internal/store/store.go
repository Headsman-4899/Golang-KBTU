package store

import (
	"context"
	"lectures/hw6/internal/models"
)

type Store interface {
	Create(ctx context.Context, game *models.Game) error
	All(ctx context.Context) ([]*models.Game, error)
	ByID(ctx context.Context, id int) (*models.Game, error)
	Update(ctx context.Context, game *models.Game) error
	Delete(ctx context.Context, id int) error
}
