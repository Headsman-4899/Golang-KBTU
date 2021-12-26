package cache

import (
	"context"
	"encoding/json"
	"lectures/hw6/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
)

type GameRedisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewGameRedisCache(host string, db int, expires time.Duration) GameCache {
	return &GameRedisCache{
		host:    host,
		db:      db,
		expires: expires,
	}
}

func (m GameRedisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     m.host,
		Password: "",
		DB:       m.db,
	})
}

func (m GameRedisCache) Set(ctx context.Context, key string, value *models.Game) {
	client := m.getClient()
	game, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	_, err = client.Set(ctx, key, game, m.expires*time.Second).Result()
	if err != nil {
		return
	}
}

func (m GameRedisCache) Get(ctx context.Context, key string) *models.Game {
	client := m.getClient()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	game := new(models.Game)
	err = json.Unmarshal([]byte(val), &game)
	if err != nil {
		panic(err)
	}
	return game
}
