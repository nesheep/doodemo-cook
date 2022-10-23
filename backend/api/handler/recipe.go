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

type recipeUsecase interface {
	Find(ctx context.Context, q string, tags []string, limit int, skip int) (entity.Recipes, int, error)
	FindOne(ctx context.Context, id string) (entity.Recipe, error)
	InsertOne(ctx context.Context, recipe entity.Recipe) (entity.Recipe, error)
	UpdateOne(ctx context.Context, id string, recipe entity.Recipe) (entity.Recipe, error)
	DeleteOne(ctx context.Context, id string) error
}

type Recipe struct {
	u recipeUsecase
}

func NewRecipe(u recipeUsecase) *Recipe {
	return &Recipe{u: u}
}

func (h *Recipe) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q, err := h.parseQeury(r.URL.Query())
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusBadRequest)
		log.Println(err)
		return
	}

	recipes, cnt, err := h.u.Find(ctx, q.q, q.tags, q.limit, q.skip())
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	res := resRecipesFromRecipes(recipes, cnt)
	response.JSON(ctx, w, res, http.StatusOK)
}

func (h *Recipe) FindOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	recipe, err := h.u.FindOne(ctx, id)
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

	res := resRecipeFromRecipe(recipe)
	response.JSON(ctx, w, res, http.StatusOK)
}

func (h *Recipe) InsertOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req reqRecipe
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromStatusCode(ctx, w, http.StatusBadRequest)
		log.Println(err)
		return
	}

	recipe := req.toRecipe()
	insertedRecipe, err := h.u.InsertOne(ctx, recipe)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	res := resRecipeFromRecipe(insertedRecipe)

	response.JSON(ctx, w, res, http.StatusCreated)
}

func (h *Recipe) UpdateOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	var req reqRecipe
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromStatusCode(ctx, w, http.StatusBadRequest)
		log.Println(err)
		return
	}

	recipe := req.toRecipe()
	updatedRecipe, err := h.u.UpdateOne(ctx, id, recipe)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	res := resRecipeFromRecipe(updatedRecipe)
	response.JSON(ctx, w, res, http.StatusOK)
}

func (h *Recipe) DeleteOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if err := h.u.DeleteOne(ctx, id); err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	res := struct {
		DeletedId string `json:"deletedId"`
	}{DeletedId: id}

	response.JSON(ctx, w, res, http.StatusOK)
}
