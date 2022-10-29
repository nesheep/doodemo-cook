package usecase

import (
	"context"
	"doodemo-cook/api/entity"
	"fmt"
)

type Recipe struct {
	r recipeRepository
}

func NewRecipe(r recipeRepository) *Recipe {
	return &Recipe{r: r}
}

func (u *Recipe) Find(ctx context.Context, q string, tags []string, limit int, skip int) (entity.Recipes, int, error) {
	recipes, err := u.r.Find(ctx, q, tags, limit, skip)
	if err != nil {
		return nil, 0, fmt.Errorf("fail 'usecase.Recipe.Find': %w", err)
	}

	cnt, err := u.r.Count(ctx, q, tags)
	if err != nil {
		return nil, 0, fmt.Errorf("fail 'usecase.Recipe.Find': %w", err)
	}

	return recipes, cnt, nil
}

func (u *Recipe) FindOne(ctx context.Context, id string) (entity.Recipe, error) {
	recipe, err := u.r.FindOne(ctx, id)
	if err != nil {
		return entity.Recipe{}, fmt.Errorf("fail 'usecase.Recipe.FindOne': %w", err)
	}

	return recipe, nil
}

func (u *Recipe) InsertOne(ctx context.Context, recipe entity.Recipe) (entity.Recipe, error) {
	recipe.Tags = recipe.Tags.Unique()
	id, err := u.r.InsertOne(ctx, recipe)
	if err != nil {
		return entity.Recipe{}, fmt.Errorf("fail 'usecase.Recipe.InsertOne': %w", err)
	}

	insertedRecipe, err := u.r.FindOne(ctx, id)
	if err != nil {
		return entity.Recipe{}, fmt.Errorf("fail 'usecase.Recipe.InsertOne': %w", err)
	}

	return insertedRecipe, nil
}

func (u *Recipe) UpdateOne(ctx context.Context, id string, recipe entity.Recipe) (entity.Recipe, error) {
	recipe.Tags = recipe.Tags.Unique()
	if err := u.r.UpdateOne(ctx, id, recipe); err != nil {
		return entity.Recipe{}, fmt.Errorf("fail 'usecase.Recipe.UpdateOne': %w", err)
	}

	updatedRecipe, err := u.r.FindOne(ctx, id)
	if err != nil {
		return entity.Recipe{}, fmt.Errorf("fail 'usecase.Recipe.UpdateOne': %w", err)
	}

	return updatedRecipe, nil
}

func (u *Recipe) DeleteOne(ctx context.Context, id string) error {
	if err := u.r.DeleteOne(ctx, id); err != nil {
		return fmt.Errorf("fail 'usecase.Recipe.DeleteOne': %w", err)
	}

	return nil
}
