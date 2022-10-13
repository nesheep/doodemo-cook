package server

import (
	"doodemo-cook/api/handler"
	"doodemo-cook/api/store"
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
		s := store.NewRecipe(db)
		h := handler.NewRecipe(s)
		r.Use(authMiddleware)
		r.Get("/", h.Find)
		r.Get("/{id}", h.FindOne)
		r.Post("/", h.InsertOne)
		r.Put("/{id}", h.UpdateOne)
		r.Delete("/{id}", h.DeleteOne)
	})

	r.Route("/tags", func(r chi.Router) {
		s := store.NewTag(db)
		h := handler.NewTag(s)
		r.Use(authMiddleware)
		r.Get("/", h.Find)
		r.Get("/{id}", h.FindOne)
		r.Post("/", h.InsertOne)
		r.Put("/{id}", h.UpdateOne)
		r.Delete("/{id}", h.DeleteOne)
	})

	return r
}
