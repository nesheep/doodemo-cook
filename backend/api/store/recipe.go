package store

import (
	"context"
	"doodemo-cook/api/entity"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const recipeColl = "recipes"

type Recipe struct {
	db *mongo.Database
}

func NewRecipe(db *mongo.Database) *Recipe {
	return &Recipe{db: db}
}

func (s *Recipe) Find(ctx context.Context, q string, limit int, skip int) (entity.RecipesWithTags, error) {
	coll := s.db.Collection(recipeColl)

	filter := bson.M{}
	pipeline := mongo.Pipeline{}

	if q != "" {
		filter = bson.M{
			"title": bson.M{"$regex": q},
		}
		pipeline = append(pipeline, bson.D{{
			Key:   "$match",
			Value: filter,
		}})
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{
			Key: "$lookup",
			Value: bson.M{
				"from":         tagColl,
				"localField":   "tags",
				"foreignField": "_id",
				"as":           "tags",
			},
		}},
	)

	total, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return entity.RecipesWithTags{}, err
	}

	cur, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return entity.RecipesWithTags{}, err
	}

	data := []entity.RecipeWithTags{}
	for cur.Next(ctx) {
		recipe := entity.NewRecipeWithTags()
		cur.Decode(&recipe)
		data = append(data, recipe)
	}

	return entity.RecipesWithTags{Data: data, Total: int(total)}, nil
}

func (s *Recipe) FindOne(ctx context.Context, id string) (entity.RecipeWithTags, error) {
	coll := s.db.Collection(recipeColl)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.RecipeWithTags{}, err
	}

	pipline := mongo.Pipeline{
		{{
			Key:   "$match",
			Value: bson.M{"_id": objId},
		}},
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         tagColl,
				"localField":   "tags",
				"foreignField": "_id",
				"as":           "tags",
			},
		}},
	}

	cur, err := coll.Aggregate(ctx, pipline)
	if err != nil {
		return entity.RecipeWithTags{}, err
	}

	recipe := entity.NewRecipeWithTags()
	if cur.Next(ctx) {
		cur.Decode(&recipe)
	} else {
		return entity.RecipeWithTags{}, fmt.Errorf("not found %s", id)
	}

	return recipe, nil
}

func (s *Recipe) InsertOne(ctx context.Context, recipe entity.Recipe) (entity.Recipe, error) {
	coll := s.db.Collection(recipeColl)

	ir, err := newInsRecipe(recipe)
	if err != nil {
		return entity.Recipe{}, err
	}

	result, err := coll.InsertOne(ctx, ir)
	if err != nil {
		return entity.Recipe{}, err
	}

	id := result.InsertedID
	newRecipe := entity.NewRecipe()
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&newRecipe); err != nil {
		return entity.Recipe{}, err
	}

	return newRecipe, nil
}

func (s *Recipe) UpdateOne(ctx context.Context, id string, recipe entity.Recipe) (entity.Recipe, error) {
	coll := s.db.Collection(recipeColl)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.Recipe{}, err
	}

	ir, err := newInsRecipe(recipe)
	if err != nil {
		return entity.Recipe{}, err
	}

	result, err := coll.ReplaceOne(ctx, bson.M{"_id": objId}, ir)
	if err != nil {
		return entity.Recipe{}, err
	}

	if result.MatchedCount < 1 {
		return entity.Recipe{}, fmt.Errorf("not found %s", id)
	}

	newRecipe := entity.NewRecipe()
	if err := coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&newRecipe); err != nil {
		return entity.Recipe{}, err
	}

	return newRecipe, nil
}

func (s *Recipe) DeleteOne(ctx context.Context, id string) (string, error) {
	coll := s.db.Collection(recipeColl)

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
