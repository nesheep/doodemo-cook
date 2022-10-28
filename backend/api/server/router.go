package server

import (
	"doodemo-cook/api/handler"
	"doodemo-cook/api/repository"
	"doodemo-cook/api/usecase"
	"doodemo-cook/lib/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(c *mongo.Client) http.Handler {
	r := chi.NewMux()
	v := validator.New()

	authMiddleware := auth.Middleware()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/tags", func(r chi.Router) {
		repo := repository.NewTag(c)
		u := usecase.NewTag(repo)
		h := handler.NewTag(u)
		r.Use(authMiddleware)
		r.Get("/", h.Find)
	})

	r.Route("/recipes", func(r chi.Router) {
		repo := repository.NewRecipe(c)
		u := usecase.NewRecipe(repo)
		h := handler.NewRecipe(u, v)
		r.Use(authMiddleware)
		r.Get("/", h.Find)
		r.Get("/{id}", h.FindOne)
		r.Post("/", h.InsertOne)
		r.Put("/{id}", h.UpdateOne)
		r.Delete("/{id}", h.DeleteOne)
	})

	return r
}
