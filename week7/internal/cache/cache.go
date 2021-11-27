package cache

import (
	"context"
	"lectures/hw6/internal/models"
)

type Cache interface {
	Close() error

	Games() GamesCacheRepo

	DeleteAll(ctx context.Context) error
}

type GamesCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.Game) error
	Get(ctx context.Context, key string) ([]*models.Game, error)
}