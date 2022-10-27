package usecase

import (
	"context"
	"doodemo-cook/api/entity"
)

type recipeRepository interface {
	Count(ctx context.Context, q string, tags []string) (int, error)
	Find(ctx context.Context, q string, tags []string, limit int, skip int) (entity.Recipes, error)
	FindOne(ctx context.Context, id string) (entity.Recipe, error)
	InsertOne(ctx context.Context, recipe entity.Recipe) (string, error)
	UpdateOne(ctx context.Context, id string, recipe entity.Recipe) error
	DeleteOne(ctx context.Context, id string) error
}

type tagRepository interface {
	Count(ctx context.Context) (int, error)
	Find(ctx context.Context) (entity.Tags, error)
}
