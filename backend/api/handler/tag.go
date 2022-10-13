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

type tagStore interface {
	Find(ctx context.Context) (entity.Tags, error)
	FindOne(ctx context.Context, id string) (entity.Tag, error)
	InsertOne(ctx context.Context, tag entity.Tag) (entity.Tag, error)
	UpdateOne(ctx context.Context, id string, recipe entity.Tag) (entity.Tag, error)
	DeleteOne(ctx context.Context, id string) (string, error)
}

type Tag struct {
	store tagStore
}

func NewTag(store tagStore) *Tag {
	return &Tag{store: store}
}

func (h *Tag) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tags, err := h.store.Find(ctx)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response.JSON(ctx, w, tags, http.StatusOK)
}

func (h *Tag) FindOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	tag, err := h.store.FindOne(ctx, id)
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

	response.JSON(ctx, w, tag, http.StatusOK)
}

func (h *Tag) InsertOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var tag entity.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		response.FromStatusCode(ctx, w, http.StatusBadRequest)
		log.Println(err)
		return
	}

	tag, err := h.store.InsertOne(ctx, tag)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response.JSON(ctx, w, tag, http.StatusCreated)
}

func (h *Tag) UpdateOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	var tag entity.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		response.FromStatusCode(ctx, w, http.StatusBadRequest)
		log.Println(err)
		return
	}

	tag, err := h.store.UpdateOne(ctx, id, tag)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response.JSON(ctx, w, tag, http.StatusOK)
}

func (h *Tag) DeleteOne(w http.ResponseWriter, r *http.Request) {
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
