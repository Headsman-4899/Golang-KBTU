package http

import (
	"context"
	"lectures/hw6/internal/cache"
	"lectures/hw6/internal/http/resources"
	"lectures/hw6/internal/store"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.GamesRepository

	cache   cache.Cache
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
	//
	//// Games
	//r.Post("/games", func(w http.ResponseWriter, r *http.Request) {
	//	game := new(models.Game)
	//	if err := json.NewDecoder(r.Body).Decode(game); err != nil {
	//		fmt.Fprintf(w, "Unknown err: %v", err)
	//		return
	//	}
	//	game.ID = primitive.NewObjectID()
	//	err := s.store.Create(r.Context(), game)
	//	if err != nil {
	//		return
	//	}
	//	w.WriteHeader(http.StatusCreated)
	//})
	//
	//r.Get("/games", func(w http.ResponseWriter, r *http.Request) {
	//	games, err := s.store.All(r.Context())
	//	if err != nil {
	//		fmt.Fprintf(w, "Unknown err: %v", err)
	//		return
	//	}
	//	render.JSON(w, r, games)
	//})
	//
	//r.Get("/games/{id}", func(w http.ResponseWriter, r *http.Request) {
	//	idStr := chi.URLParam(r, "id")
	//
	//	game, error := s.store.ByID(r.Context(), idStr)
	//	if error != nil {
	//		fmt.Fprintf(w, "Unknown err: %v", error)
	//		return
	//	}
	//	render.JSON(w, r, game)
	//})
	//
	//r.Put("/games", func(w http.ResponseWriter, r *http.Request) {
	//	game := new(models.Game)
	//	if err := json.NewDecoder(r.Body).Decode(game); err != nil {
	//		fmt.Fprintf(w, "Unknown err: %v", err)
	//		return
	//	}
	//
	//	s.store.Update(r.Context(), game)
	//})
	//
	//r.Delete("/games/{id}", func(w http.ResponseWriter, r *http.Request) {
	//	idStr := chi.URLParam(r, "id")
	//	//id, err := strconv.Atoi(idStr)
	//
	//	error := s.store.Delete(r.Context(), idStr)
	//	if error != nil {
	//		return
	//	}
	//})

	// Profiles
	//r.Post("/profiles", func(w http.ResponseWriter, r *http.Request) {
	//	profile := new(models.Profile)
	//	if err := json.NewDecoder(r.Body).Decode(profile); err != nil {
	//		fmt.Fprintf(w, "Unknown err: %v", err)
	//		return
	//	}
	//
	//	s.store.Profiles().Create(r.Context(), profile)
	//})
	//r.Get("/profiles", func(w http.ResponseWriter, r *http.Request) {
	//	profiles, err := s.store.Profiles().All(r.Context())
	//	if err != nil {
	//		fmt.Fprintf(w, "Unknown err: %v", err)
	//		return
	//	}
	//
	//	render.JSON(w, r, profiles)
	//})
	//r.Get("/profiles/{id}", func(w http.ResponseWriter, r *http.Request) {
	//	idStr := chi.URLParam(r, "id")
	//	id, err := strconv.Atoi(idStr)
	//	if err != nil {
	//		fmt.Fprintf(w, "Unknown err: %v", err)
	//		return
	//	}
	//	profile, err := s.store.Profiles().ByID(r.Context(), id)
	//	if err != nil {
	//		fmt.Fprintf(w, "Unknown err: %v", err)
	//		return
	//	}
	//
	//	render.JSON(w, r, profile)
	//})
	//r.Put("/profiles", func(w http.ResponseWriter, r *http.Request) {
	//	profile := new(models.Profile)
	//	if err := json.NewDecoder(r.Body).Decode(profile); err != nil {
	//		fmt.Fprintf(w, "Unknown err: %v", err)
	//		return
	//	}
	//
	//	s.store.Profiles().Update(r.Context(), profile)
	//})
	//r.Delete("/profiles/{id}", func(w http.ResponseWriter, r *http.Request) {
	//	idStr := chi.URLParam(r, "id")
	//	id, err := strconv.Atoi(idStr)
	//
	//	if err != nil {
	//		fmt.Fprintf(w, "Unknown err: %v", err)
	//		return
	//	}
	//
	//	error := s.store.Profiles().Delete(r.Context(), id)
	//	if error != nil {
	//		return
	//	}
	//})
	gamesResource := resources.NewGamesResource(s.store, s.cache)
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
	os.RemoveAll("./tmp")
}
