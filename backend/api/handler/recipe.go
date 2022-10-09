package handler

import (
	"context"
	"doodemo-cook/api/entity"
	"doodemo-cook/lib/response"
	"encoding/json"
	"log"
	"net/http"
)

type recipeStore interface {
	Find(ctx context.Context) ([]entity.Recipe, error)
	InsertOne(ctx context.Context, recipe entity.Recipe) (entity.Recipe, error)
}

type Recipe struct {
	store recipeStore
}

func NewRecipe(store recipeStore) *Recipe {
	return &Recipe{store: store}
}

func (h *Recipe) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rs, err := h.store.Find(ctx)
	if err != nil {
		response.JSON(ctx, w, response.ErrorResponse{Message: "internal server error"}, http.StatusInternalServerError)
		log.Panicln(err)
		return
	}

	recipes := entity.Recipes{
		Data:  rs,
		Total: len(rs),
	}

	response.JSON(ctx, w, recipes, http.StatusOK)
}

func (h *Recipe) InsertOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var recipe entity.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		response.JSON(ctx, w, response.ErrorResponse{Message: "bad request"}, http.StatusBadRequest)
		log.Panicln(err)
		return
	}

	recipe, err := h.store.InsertOne(ctx, recipe)
	if err != nil {
		response.JSON(ctx, w, response.ErrorResponse{Message: "internal server error"}, http.StatusInternalServerError)
		log.Panicln(err)
		return
	}

	response.JSON(ctx, w, recipe, http.StatusCreated)
}
