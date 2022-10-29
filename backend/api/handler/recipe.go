package handler

import (
	"doodemo-cook/lib/response"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type Recipe struct {
	u recipeUsecase
	v *validator.Validate
}

func NewRecipe(u recipeUsecase, v *validator.Validate) *Recipe {
	return &Recipe{u: u, v: v}
}

func (h *Recipe) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q, err := h.parseQuery(r.URL.Query())
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusBadRequest)
		log.Printf("fail 'handler.Recipe.Find': %v", err)
		return
	}

	recipes, cnt, err := h.u.Find(ctx, q.q, q.tags, q.limit, q.skip())
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Printf("fail 'handler.Recipe.Find': %v", err)
		return
	}

	res := resRecipesFromRecipes(recipes, cnt)
	response.JSON(ctx, w, res, http.StatusOK)
}

func (h *Recipe) FindOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	recipe, err := h.u.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			response.FromStatusCode(ctx, w, http.StatusNotFound)
			log.Printf("fail 'handler.Recipe.FindOne': %v", err)
			return
		}
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Printf("fail 'handler.Recipe.FindOne': %v", err)
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
		log.Printf("fail 'handler.Recipe.InsertOne': %v", err)
		return
	}

	if err := h.v.Struct(req); err != nil {
		response.FromStatusCode(ctx, w, http.StatusBadRequest)
		log.Printf("fail 'handler.Recipe.InsertOne': %v", err)
		return
	}

	recipe := req.toRecipe()
	insertedRecipe, err := h.u.InsertOne(ctx, recipe)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Printf("fail 'handler.Recipe.InsertOne': %v", err)
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
		log.Printf("fail 'handler.Recipe.UpdateOne': %v", err)
		return
	}

	if err := h.v.Struct(req); err != nil {
		response.FromStatusCode(ctx, w, http.StatusBadRequest)
		log.Printf("fail 'handler.Recipe.UpdateOne': %v", err)
		return
	}

	recipe := req.toRecipe()
	updatedRecipe, err := h.u.UpdateOne(ctx, id, recipe)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Printf("fail 'handler.Recipe.UpdateOne': %v", err)
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
		log.Printf("fail 'handler.Recipe.DeleteOne': %v", err)
		return
	}

	res := struct {
		DeletedId string `json:"deletedId"`
	}{DeletedId: id}

	response.JSON(ctx, w, res, http.StatusOK)
}
