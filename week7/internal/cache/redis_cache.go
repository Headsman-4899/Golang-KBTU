package cache

import (
	"context"
	"lectures/hw6/internal/models"
)

type Cache interface {
	Games() GameCache
}

type GameCache interface {
	Set(ctx context.Context, key string, value *models.Game)
	Get(ctx context.Context, key string) *models.Game
}
