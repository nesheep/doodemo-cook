package repository

import (
	"context"
	"doodemo-cook/api/entity"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const recipeColl = "recipes"

var lookupStage = bson.D{{
	Key: "$lookup",
	Value: bson.M{
		"from":         tagColl,
		"localField":   "tags",
		"foreignField": "_id",
		"as":           "tags",
	},
}}

type Recipe struct {
	db  *mongo.Database
	tag *Tag
}

func NewRecipe(db *mongo.Database, tag *Tag) *Recipe {
	return &Recipe{db: db, tag: tag}
}

func (r *Recipe) buildFilter(q string, tags []string) (bson.M, error) {
	filter := bson.M{}

	if q != "" {
		filter["title"] = bson.M{"$regex": q}
	}

	if len(tags) > 0 {
		t := bson.A{}
		for _, v := range tags {
			oid, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				return bson.M{}, err
			}
			t = append(t, oid)
		}
		filter["tags"] = bson.M{"$all": t}
	}

	return filter, nil
}

func (r *Recipe) Count(ctx context.Context, q string, tags []string) (int, error) {
	coll := r.db.Collection(recipeColl)

	filter, err := r.buildFilter(q, tags)
	if err != nil {
		return 0, err
	}

	cnt, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(cnt), nil
}

func (r *Recipe) Find(ctx context.Context, q string, tags []string, limit int, skip int) (entity.Recipes, error) {
	coll := r.db.Collection(recipeColl)

	filter, err := r.buildFilter(q, tags)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{}

	if len(filter) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: filter}})
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
		lookupStage,
	)

	recipes := entity.Recipes{}
	cur, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var b bRecipe
		cur.Decode(&b)
		recipes = append(recipes, b.toRecipe())
	}

	return recipes, nil
}

func (r *Recipe) FindOne(ctx context.Context, id string) (entity.Recipe, error) {
	coll := r.db.Collection(recipeColl)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.Recipe{}, err
	}

	pipline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": oid}}},
		lookupStage,
	}

	cur, err := coll.Aggregate(ctx, pipline)
	if err != nil {
		return entity.Recipe{}, err
	}

	var b bRecipe
	if cur.Next(ctx) {
		cur.Decode(&b)
	} else {
		return entity.Recipe{}, fmt.Errorf("not found %s", id)
	}

	return b.toRecipe(), nil
}

func (r *Recipe) InsertOne(ctx context.Context, recipe entity.Recipe) (string, error) {
	coll := r.db.Collection(recipeColl)

	for i, v := range recipe.Tags {
		id, err := r.tag.insertIfNotExists(ctx, v)
		if err != nil {
			return "", err
		}
		recipe.Tags[i].ID = id
	}

	b, err := bInputRecipeFromRecipe(recipe)
	if err != nil {
		return "", err
	}

	result, err := coll.InsertOne(ctx, b)
	if err != nil {
		return "", err
	}

	id := result.InsertedID
	oid, ok := id.(primitive.ObjectID)
	if !ok {
		return "", errors.New("read inserted ID error")
	}

	return oid.Hex(), nil
}

func (r *Recipe) UpdateOne(ctx context.Context, id string, recipe entity.Recipe) error {
	coll := r.db.Collection(recipeColl)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	for i, v := range recipe.Tags {
		id, err := r.tag.insertIfNotExists(ctx, v)
		if err != nil {
			return err
		}
		recipe.Tags[i].ID = id
	}

	b, err := bInputRecipeFromRecipe(recipe)
	if err != nil {
		return err
	}

	result, err := coll.ReplaceOne(ctx, bson.M{"_id": oid}, b)
	if err != nil {
		return err
	}

	if result.MatchedCount < 1 {
		return fmt.Errorf("not found %s", id)
	}

	return nil
}

func (r *Recipe) DeleteOne(ctx context.Context, id string) error {
	coll := r.db.Collection(recipeColl)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if result.DeletedCount < 1 {
		return fmt.Errorf("not found %s", id)
	}

	return nil
}
