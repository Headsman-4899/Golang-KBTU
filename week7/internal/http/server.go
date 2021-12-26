package http

import (
	"context"
	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
	"lectures/hw6/internal/http/resources"
	"lectures/hw6/internal/message_broker"
	"lectures/hw6/internal/store"
	"log"
	"net/http"
	"time"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.GamesRepository

	cache   *lru.TwoQueueCache
	broker  message_broker.MessageBroker
	Address string
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()
	gamesResource := resources.NewGamesResource(s.store, s.broker, s.cache)
	r.Mount("/games", gamesResource.Routes())

	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}

	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
	//os.RemoveAll("./tmp")
}
