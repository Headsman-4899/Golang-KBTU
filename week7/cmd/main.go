package main

import (
	"context"
	redis_cache "lectures/hw6/internal/cache/redis-cache"
	"lectures/hw6/internal/http"
	"lectures/hw6/internal/store/inmemory"
	"log"
)

func main() {
	store := inmemory.Init()

	//srv := http.NewServer(context.Background(), ":8080", store)
	//if err := srv.Run(); err != nil {
	//	log.Println(err)
	//}

	cache := redis_cache.NewRedisCache(":8080", 1, 1800)

	srv := http.NewServer(context.Background(),
		http.WithAddress(":8080"),
		http.WithStore(store),
		http.WithCache(cache),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}
