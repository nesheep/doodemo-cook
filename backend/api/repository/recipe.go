package repository

import (
	"context"
	"doodemo-cook/api/entity"

	"go.mongodb.org/mongo-driver/mongo"
)

type Recipe struct {
	db *mongo.Database
}

func NewRecipe(db *mongo.Database) *Recipe {
	return &Recipe{db: db}
}

func (r *Recipe) Count(ctx context.Context, q string, tags []string) (int, error) {
	return 0, nil
}

func (r *Recipe) Find(ctx context.Context, q string, tags []string, limit int, skip int) (entity.Recipes, error) {
	recipes := entity.Recipes{}
	return recipes, nil
}

func (r *Recipe) FindOne(ctx context.Context, id string) (entity.Recipe, error) {
	recipe := entity.Recipe{}
	return recipe, nil
}

func (r *Recipe) InsertOne(ctx context.Context, recipe entity.Recipe) (string, error) {
	return "", nil
}

func (r *Recipe) UpdateOne(ctx context.Context, id string, recipe entity.Recipe) error {
	return nil
}

func (r *Recipe) DeleteOne(ctx context.Context, id string) error {
	return nil
}
