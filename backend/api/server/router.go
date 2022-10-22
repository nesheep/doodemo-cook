package server

import (
	"doodemo-cook/api/handler"
	"doodemo-cook/api/repository"
	"doodemo-cook/api/usecase"
	"doodemo-cook/lib/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(db *mongo.Database) http.Handler {
	r := chi.NewMux()

	authMiddleware := auth.Middleware()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/recipes", func(r chi.Router) {
		repo := repository.NewRecipe(db)
		u := usecase.NewRecipe(repo)
		h := handler.NewRecipe(u)
		r.Use(authMiddleware)
		r.Get("/", h.Find)
		r.Get("/{id}", h.FindOne)
		r.Post("/", h.InsertOne)
		r.Put("/{id}", h.UpdateOne)
		r.Delete("/{id}", h.DeleteOne)
	})

	r.Route("/tags", func(r chi.Router) {
		repo := repository.NewTag(db)
		u := usecase.NewTag(repo)
		h := handler.NewTag(u)
		r.Use(authMiddleware)
		r.Get("/", h.Find)
	})

	return r
}
