package redis_cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"lectures/hw6/internal/cache"
	"lectures/hw6/internal/models"
	"time"
)

func (rc RedisCache) Games() cache.GamesCacheRepo {
	if rc.games == nil {
		rc.games = newGamesRepo(rc.client, rc.expires)
	}

	return rc.games
}

type GamesRepo struct {
	client  *redis.Client
	expires time.Duration
}

func newGamesRepo(client *redis.Client, exp time.Duration) cache.GamesCacheRepo {
	return &GamesRepo{
		client:  client,
		expires: exp,
	}
}

func (c GamesRepo) Set(ctx context.Context, key string, value []*models.Game) error {
	gamesBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = c.client.Set(ctx, key, gamesBytes, c.expires*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func (c GamesRepo) Get(ctx context.Context, key string) ([]*models.Game, error) {
	result, err := c.client.Get(ctx, key).Result()
	switch err {
	case nil:
		break
	case redis.Nil:
		return nil, nil
	default:
		return nil, err
	}

	games := make([]*models.Game, 0)
	if err = json.Unmarshal([]byte(result), &games); err != nil {
		return nil, err
	}

	return games, nil
}
