package store

import (
	"context"
	"doodemo-cook/api/entity"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *Recipe) FindOne(ctx context.Context, id string) (entity.Recipe, error) {
	coll := s.db.Collection(name)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.Recipe{}, err
	}

	var newRecipe entity.Recipe
	if err := coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&newRecipe); err != nil {
		return entity.Recipe{}, err
	}

	return newRecipe, nil
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

func (s *Recipe) UpdateOne(ctx context.Context, id string, recipe entity.Recipe) (entity.Recipe, error) {
	coll := s.db.Collection(name)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.Recipe{}, err
	}

	result, err := coll.ReplaceOne(ctx, bson.M{"_id": objId}, recipe)
	if err != nil {
		return entity.Recipe{}, err
	}

	if result.MatchedCount < 1 {
		return entity.Recipe{}, fmt.Errorf("not found %s", id)
	}

	var newRecipe entity.Recipe
	if err := coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&newRecipe); err != nil {
		return entity.Recipe{}, err
	}

	return newRecipe, nil
}

func (s *Recipe) DeleteOne(ctx context.Context, id string) (string, error) {
	coll := s.db.Collection(name)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	result, err := coll.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return "", err
	}

	if result.DeletedCount < 1 {
		return "", fmt.Errorf("not found %s", id)
	}

	return id, nil
}
