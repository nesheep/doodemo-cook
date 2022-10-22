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

type Recipe struct {
	r recipeRepository
}

func NewRecipe(r recipeRepository) *Recipe {
	return &Recipe{r: r}
}

func (u *Recipe) Find(ctx context.Context, q string, tags []string, limit int, skip int) (entity.Recipes, int, error) {
	recipes, err := u.r.Find(ctx, q, tags, limit, skip)
	if err != nil {
		return entity.Recipes{}, 0, err
	}

	cnt, err := u.r.Count(ctx, q, tags)
	if err != nil {
		return entity.Recipes{}, 0, err
	}

	return recipes, cnt, nil
}

func (u *Recipe) FindOne(ctx context.Context, id string) (entity.Recipe, error) {
	recipe, err := u.r.FindOne(ctx, id)
	if err != nil {
		return entity.Recipe{}, err
	}

	return recipe, nil
}

func (u *Recipe) InsertOne(ctx context.Context, recipe entity.Recipe) (entity.Recipe, error) {
	id, err := u.r.InsertOne(ctx, recipe)
	if err != nil {
		return entity.Recipe{}, err
	}

	insertedRecipe, err := u.r.FindOne(ctx, id)
	if err != nil {
		return entity.Recipe{}, err
	}

	return insertedRecipe, nil
}

func (u *Recipe) UpdateOne(ctx context.Context, id string, recipe entity.Recipe) (entity.Recipe, error) {
	if err := u.r.UpdateOne(ctx, id, recipe); err != nil {
		return entity.Recipe{}, err
	}

	updatedRecipe, err := u.r.FindOne(ctx, id)
	if err != nil {
		return entity.Recipe{}, err
	}

	return updatedRecipe, nil
}

func (u *Recipe) DeleteOne(ctx context.Context, id string) error {
	if err := u.r.DeleteOne(ctx, id); err != nil {
		return err
	}

	return nil
}
