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

	tr := repository.NewTag(db)
	rr := repository.NewRecipe(db, tr)

	r.Route("/tags", func(r chi.Router) {
		u := usecase.NewTag(tr)
		h := handler.NewTag(u)
		r.Use(authMiddleware)
		r.Get("/", h.Find)
	})

	r.Route("/recipes", func(r chi.Router) {
		u := usecase.NewRecipe(rr)
		h := handler.NewRecipe(u)
		r.Use(authMiddleware)
		r.Get("/", h.Find)
		r.Get("/{id}", h.FindOne)
		r.Post("/", h.InsertOne)
		r.Put("/{id}", h.UpdateOne)
		r.Delete("/{id}", h.DeleteOne)
	})

	return r
}
