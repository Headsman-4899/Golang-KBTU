package resources

import (
	"encoding/json"
	"fmt"
	"lectures/hw6/internal/message_broker"
	"lectures/hw6/internal/models"
	"lectures/hw6/internal/store"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	lru "github.com/hashicorp/golang-lru"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GamesResource struct {
	store  store.GamesRepository
	broker message_broker.MessageBroker
	cache  *lru.TwoQueueCache
}

func NewGamesResource(store store.GamesRepository, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *GamesResource {
	return &GamesResource{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (cr *GamesResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", cr.CreateGame)
	r.Get("/", cr.AllGames)
	r.Get("/{id}", cr.ById)
	r.Put("/", cr.UpdateGame)
	r.Delete("/{id}", cr.DeleteGame)

	return r
}

func (cr *GamesResource) CreateGame(w http.ResponseWriter, r *http.Request) {
	game := new(models.Game)
	game.ID = primitive.NewObjectID()

	if err := json.NewDecoder(r.Body).Decode(game); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Create(r.Context(), game); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	err := cr.broker.Cache().Purge()
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (cr *GamesResource) AllGames(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.GamesFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		gamesFromCache, ok := cr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, gamesFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	games, err := cr.store.All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" && len(games) > 0 {
		// err = cr.cache.Games().Set(r.Context(), searchQuery, games)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	fmt.Fprintf(w, "Cache err: %v", err)
		// 	return
		// }
		cr.cache.Add(searchQuery, games)
	}

	render.JSON(w, r, games)
}

func (cr *GamesResource) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	game, err := cr.store.ByID(r.Context(), idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	render.JSON(w, r, game)
}

func (cr *GamesResource) UpdateGame(w http.ResponseWriter, r *http.Request) {
	game := new(models.Game)
	if err := json.NewDecoder(r.Body).Decode(game); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(game,
		validation.Field(&game.ID, validation.Required),
		validation.Field(&game.Name, validation.Required))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = cr.store.Update(r.Context(), game); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	err = cr.broker.Cache().Remove(game.ID)
	if err != nil {
		return
	}
}

func (cr *GamesResource) DeleteGame(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	//id, err := strconv.Atoi(idStr)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	fmt.Fprintf(w, "Unknown err: %v", err)
	//	return
	//}

	if err := cr.store.Delete(r.Context(), idStr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	err := cr.broker.Cache().Remove(idStr)
	if err != nil {
		return
	}
}
