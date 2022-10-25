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

func NewRecipe(db *mongo.Database) *Recipe {
	return &Recipe{db: db, tag: &Tag{db: db}}
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
	var id any
	err := r.db.Client().UseSession(ctx, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return err
		}

		coll := r.db.Collection(recipeColl)

		for i, v := range recipe.Tags {
			id, err := r.tag.insertOrIncrement(sc, v)
			if err != nil {
				sc.AbortTransaction(context.Background())
				return err
			}
			recipe.Tags[i].ID = id
		}

		b, err := bInputRecipeFromRecipe(recipe)
		if err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		result, err := coll.InsertOne(sc, b)
		if err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		id = result.InsertedID

		sc.CommitTransaction(context.Background())
		return nil
	})
	if err != nil {
		return "", err
	}

	oid, ok := id.(primitive.ObjectID)
	if !ok {
		return "", errors.New("read inserted ID error")
	}

	return oid.Hex(), nil
}

func (r *Recipe) UpdateOne(ctx context.Context, id string, recipe entity.Recipe) error {
	err := r.db.Client().UseSession(ctx, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return err
		}

		coll := r.db.Collection(recipeColl)

		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		old, err := r.FindOne(sc, id)
		if err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		for i, v := range recipe.Tags {
			exists := false
			for _, ov := range old.Tags {
				if v.Name == ov.Name {
					exists = true
					break
				}
			}
			if exists {
				tag, err := r.tag.findByName(sc, v.Name)
				if err != nil {
					sc.AbortTransaction(context.Background())
					return err
				}
				recipe.Tags[i].ID = tag.ID
				continue
			}

			tagId, err := r.tag.insertOrIncrement(sc, v)
			if err != nil {
				sc.AbortTransaction(context.Background())
				return err
			}
			recipe.Tags[i].ID = tagId
		}

		for _, ov := range old.Tags {
			exists := false
			for _, v := range recipe.Tags {
				if ov.Name == v.Name {
					exists = true
					break
				}
			}
			if exists {
				continue
			}

			if err := r.tag.deleteOrDecrement(sc, ov.ID); err != nil {
				sc.AbortTransaction(context.Background())
				return err
			}
		}

		b, err := bInputRecipeFromRecipe(recipe)
		if err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		result, err := coll.ReplaceOne(sc, bson.M{"_id": oid}, b)
		if err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		if result.MatchedCount < 1 {
			sc.AbortTransaction(context.Background())
			return err
		}

		sc.CommitTransaction(context.Background())
		return nil
	})
	if err != nil {
		return fmt.Errorf("err repository.Recipe.UpdateOne: %v", err)
	}

	return nil
}

func (r *Recipe) DeleteOne(ctx context.Context, id string) error {
	err := r.db.Client().UseSession(ctx, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return err
		}

		coll := r.db.Collection(recipeColl)

		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		recipe, err := r.FindOne(sc, id)
		if err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		for _, v := range recipe.Tags {
			if err := r.tag.deleteOrDecrement(sc, v.ID); err != nil {
				sc.AbortTransaction(context.Background())
				return err
			}
		}

		result, err := coll.DeleteOne(sc, bson.M{"_id": oid})
		if err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		if result.DeletedCount < 1 {
			sc.AbortTransaction(context.Background())
			return fmt.Errorf("not found %s", id)
		}

		sc.CommitTransaction(context.Background())
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
