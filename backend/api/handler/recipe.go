package handler

import (
	"context"
	"doodemo-cook/api/entity"
	"doodemo-cook/lib/response"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type recipeStore interface {
	Find(ctx context.Context) (entity.RecipesWithTags, error)
	FindOne(ctx context.Context, id string) (entity.RecipeWithTags, error)
	InsertOne(ctx context.Context, recipe entity.Recipe) (entity.Recipe, error)
	UpdateOne(ctx context.Context, id string, recipe entity.Recipe) (entity.Recipe, error)
	DeleteOne(ctx context.Context, id string) (string, error)
}

type Recipe struct {
	store recipeStore
}

func NewRecipe(store recipeStore) *Recipe {
	return &Recipe{store: store}
}

func (h *Recipe) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	recipes, err := h.store.Find(ctx)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response.JSON(ctx, w, recipes, http.StatusOK)
}

func (h *Recipe) FindOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	recipe, err := h.store.FindOne(ctx, id)
	if err == mongo.ErrNoDocuments {
		response.FromStatusCode(ctx, w, http.StatusNotFound)
		log.Println(err)
		return
	}
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response.JSON(ctx, w, recipe, http.StatusOK)
}

func (h *Recipe) InsertOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	recipe := entity.NewRecipe()
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		response.FromStatusCode(ctx, w, http.StatusBadRequest)
		log.Println(err)
		return
	}

	recipe, err := h.store.InsertOne(ctx, recipe)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response.JSON(ctx, w, recipe, http.StatusCreated)
}

func (h *Recipe) UpdateOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	recipe := entity.NewRecipe()
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		response.FromStatusCode(ctx, w, http.StatusBadRequest)
		log.Println(err)
		return
	}

	recipe, err := h.store.UpdateOne(ctx, id, recipe)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response.JSON(ctx, w, recipe, http.StatusOK)
}

func (h *Recipe) DeleteOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	deletedId, err := h.store.DeleteOne(ctx, id)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	b := struct {
		DeletedId string `json:"deletedId"`
	}{DeletedId: deletedId}

	response.JSON(ctx, w, b, http.StatusOK)
}
