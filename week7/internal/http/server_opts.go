package http

import (
	//"lectures/hw6/internal/cache"
	"lectures/hw6/internal/store"

	lru "github.com/hashicorp/golang-lru"
)

type ServerOption func(srv *Server)

func WithAddress(address string) ServerOption {
	return func(srv *Server) {
		srv.Address = address
	}
}

func WithStore(store store.GamesRepository) ServerOption {
	return func(srv *Server) {
		srv.store = store
	}
}

func WithCache(cache *lru.TwoQueueCache) ServerOption {
	return func(srv *Server) {
		srv.cache = cache
	}
}
