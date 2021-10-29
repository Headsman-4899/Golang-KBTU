package http

import (
	"context"
	"encoding/json"
	"fmt"
	"lectures/hw6/internal/models"
	"lectures/hw6/internal/store"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store

	Address string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
		store:       store,

		Address: address,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()


	// Games
	r.Post("/games", func(w http.ResponseWriter, r *http.Request) {
		game := new(models.Game)
		if err := json.NewDecoder(r.Body).Decode(game); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Games().Create(r.Context(), game)
	})
	r.Get("/games", func(w http.ResponseWriter, r *http.Request) {
		games, err := s.store.Games().All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, games)
	})
	r.Get("/games/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		game, err := s.store.Games().ByID(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		render.JSON(w, r, game)
	})
	r.Put("/games", func(w http.ResponseWriter, r *http.Request) {
		game := new(models.Game)
		if err := json.NewDecoder(r.Body).Decode(game); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Games().Update(r.Context(), game)
	})
	r.Delete("/games/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Games().Delete(r.Context(), id)
	})



	// Profiles
	r.Post("/profiles", func(w http.ResponseWriter, r *http.Request) {
		profile := new(models.Profile)
		if err := json.NewDecoder(r.Body).Decode(profile); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Profiles().Create(r.Context(), profile)
	})
	r.Get("/profiles", func(w http.ResponseWriter, r *http.Request) {
		profiles, err := s.store.Profiles().All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, profiles)
	})
	r.Get("/profiles/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		profile, err := s.store.Profiles().ByID(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, profile)
	})
	r.Put("/profiles", func(w http.ResponseWriter, r *http.Request) {
		profile := new(models.Profile)
		if err := json.NewDecoder(r.Body).Decode(profile); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Profiles().Update(r.Context(), profile)
	})
	r.Delete("/profiles/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Profiles().Delete(r.Context(), id)
	})

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
}
