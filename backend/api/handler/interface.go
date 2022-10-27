package handler

import (
	"context"
	"doodemo-cook/api/entity"
)

type recipeUsecase interface {
	Find(ctx context.Context, q string, tags []string, limit int, skip int) (entity.Recipes, int, error)
	FindOne(ctx context.Context, id string) (entity.Recipe, error)
	InsertOne(ctx context.Context, recipe entity.Recipe) (entity.Recipe, error)
	UpdateOne(ctx context.Context, id string, recipe entity.Recipe) (entity.Recipe, error)
	DeleteOne(ctx context.Context, id string) error
}

type tagUsecase interface {
	Find(ctx context.Context) (entity.Tags, int, error)
}
