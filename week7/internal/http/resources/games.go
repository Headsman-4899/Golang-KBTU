package resources

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lectures/hw6/internal/cache"
	"lectures/hw6/internal/models"
	"lectures/hw6/internal/store"
	"net/http"
)

type GamesResource struct {
	store store.GamesRepository
	cache cache.Cache
}

func NewGamesResource(store store.GamesRepository, cache cache.Cache) *GamesResource {
	return &GamesResource{
		store: store,
		cache: cache,
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

	if err := cr.cache.DeleteAll(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Cache err: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cr *GamesResource) AllGames(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.GamesFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		gamesFromCache, err := cr.cache.Games().Get(r.Context(), searchQuery)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Cache err: %v", err)
			return
		}
		if gamesFromCache != nil {
			render.JSON(w, r, gamesFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	categories, err := cr.store.All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" && len(categories) > 0 {
		err = cr.cache.Games().Set(r.Context(), searchQuery, categories)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Cache err: %v", err)
			return
		}
	}

	render.JSON(w, r, categories)
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
	category := new(models.Game)
	if err := json.NewDecoder(r.Body).Decode(category); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(category,
		validation.Field(&category.ID, validation.Required),
		validation.Field(&category.Name, validation.Required))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = cr.store.Update(r.Context(), category); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
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
}