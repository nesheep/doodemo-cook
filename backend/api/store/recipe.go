package store

import (
	"context"
	"doodemo-cook/api/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const name = "recipes"

type Recipe struct {
	db *mongo.Database
}

func NewRecipe(db *mongo.Database) *Recipe {
	return &Recipe{db: db}
}

func (s *Recipe) Find(ctx context.Context) (entity.Recipes, error) {
	coll := s.db.Collection(name)

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return entity.Recipes{}, err
	}

	data := []entity.Recipe{}
	for cur.Next(ctx) {
		var recipe entity.Recipe
		cur.Decode(&recipe)
		data = append(data, recipe)
	}

	return entity.Recipes{Data: data, Total: len(data)}, nil
}

func (s *Recipe) InsertOne(ctx context.Context, recipe entity.Recipe) (entity.Recipe, error) {
	coll := s.db.Collection(name)

	result, err := coll.InsertOne(ctx, recipe)
	if err != nil {
		return entity.Recipe{}, err
	}

	id := result.InsertedID
	var newRecipe entity.Recipe
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&newRecipe); err != nil {
		return entity.Recipe{}, err
	}

	return newRecipe, nil
}
