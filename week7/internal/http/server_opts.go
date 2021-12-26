package http

import (
	lru "github.com/hashicorp/golang-lru"
	"lectures/hw6/internal/message_broker"
	"lectures/hw6/internal/store"
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

func WithBroker(broker message_broker.MessageBroker) ServerOption {
	return func(srv *Server) {
		srv.broker = broker
	}
}
