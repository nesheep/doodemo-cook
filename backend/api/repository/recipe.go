package repository

import (
	"context"
	"doodemo-cook/api/entity"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var tagLookupStage = bson.D{{
	Key: "$lookup",
	Value: bson.M{
		"from":         tagColl,
		"localField":   "tags",
		"foreignField": "_id",
		"as":           "tags",
	},
}}

type Recipe struct {
	c   *mongo.Client
	tag *Tag
}

func NewRecipe(c *mongo.Client) *Recipe {
	return &Recipe{c: c, tag: &Tag{c: c}}
}

func (r *Recipe) coll() *mongo.Collection {
	return r.c.Database(dbName).Collection(recipeColl)
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
				return bson.M{}, fmt.Errorf("fail 'repository.Recipe.buildFilter': %w", err)
			}
			t = append(t, oid)
		}
		filter["tags"] = bson.M{"$all": t}
	}

	return filter, nil
}

func (r *Recipe) Count(ctx context.Context, q string, tags []string) (int, error) {
	coll := r.coll()

	filter, err := r.buildFilter(q, tags)
	if err != nil {
		return 0, fmt.Errorf("fail 'repository.Recipe.Count': %w", err)
	}

	cnt, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("fail 'repository.Recipe.Count': %w", err)
	}

	return int(cnt), nil
}

func (r *Recipe) Find(ctx context.Context, q string, tags []string, limit int, skip int) (entity.Recipes, error) {
	coll := r.coll()

	filter, err := r.buildFilter(q, tags)
	if err != nil {
		return nil, fmt.Errorf("fail 'repository.Recipe.Find': %w", err)
	}

	pipeline := mongo.Pipeline{}

	if len(filter) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: filter}})
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
		tagLookupStage,
	)

	recipes := entity.Recipes{}
	cur, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("fail 'repository.Recipe.Find': %w", err)
	}

	for cur.Next(ctx) {
		var b bRecipe
		cur.Decode(&b)
		recipes = append(recipes, b.toRecipe())
	}

	return recipes, nil
}

func (r *Recipe) FindOne(ctx context.Context, id string) (entity.Recipe, error) {
	coll := r.coll()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.Recipe{}, fmt.Errorf("fail 'repository.Recipe.FindOne': %w", err)
	}

	pipline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": oid}}},
		tagLookupStage,
	}

	cur, err := coll.Aggregate(ctx, pipline)
	if err != nil {
		return entity.Recipe{}, fmt.Errorf("fail 'repository.Recipe.FindOne': %w", err)
	}

	var b bRecipe
	if cur.Next(ctx) {
		cur.Decode(&b)
	} else {
		return entity.Recipe{}, fmt.Errorf("fail 'repository.Recipe.FindOne': %w", mongo.ErrNoDocuments)
	}

	return b.toRecipe(), nil
}

func (r *Recipe) InsertOne(ctx context.Context, recipe entity.Recipe) (string, error) {
	var id any
	err := r.c.UseSession(ctx, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return err
		}

		coll := r.coll()

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
		return "", fmt.Errorf("fail 'repository.Recipe.InsertOne': %w", err)
	}

	oid, ok := id.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("fail 'repository.Recipe.InsertOne': %w", err)
	}

	return oid.Hex(), nil
}

func (r *Recipe) UpdateOne(ctx context.Context, id string, recipe entity.Recipe) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("fail 'repository.Recipe.UpdateOne': %w", err)
	}

	err = r.c.UseSession(ctx, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return err
		}

		coll := r.coll()

		old, err := r.FindOne(sc, id)
		if err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		for i, v := range recipe.Tags {
			exists := false
			for _, ov := range old.Tags {
				if v.Name == ov.Name {
					recipe.Tags[i].ID = ov.ID
					exists = true
					break
				}
			}

			if !exists {
				tagId, err := r.tag.insertOrIncrement(sc, v)
				if err != nil {
					sc.AbortTransaction(context.Background())
					return err
				}
				recipe.Tags[i].ID = tagId
			}
		}

		for _, ov := range old.Tags {
			exists := false
			for _, v := range recipe.Tags {
				if ov.Name == v.Name {
					exists = true
					break
				}
			}

			if !exists {
				if err := r.tag.deleteOrDecrement(sc, ov.ID); err != nil {
					sc.AbortTransaction(context.Background())
					return err
				}
			}
		}

		b, err := bInputRecipeFromRecipe(recipe)
		if err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		if _, err = coll.ReplaceOne(sc, bson.M{"_id": oid}, b); err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		sc.CommitTransaction(context.Background())
		return nil
	})
	if err != nil {
		return fmt.Errorf("fail 'repository.Recipe.UpdateOne': %w", err)
	}

	return nil
}

func (r *Recipe) DeleteOne(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("fail 'repository.Recipe.DeleteOne': %w", err)
	}

	err = r.c.UseSession(ctx, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return err
		}

		coll := r.coll()

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

		if _, err := coll.DeleteOne(sc, bson.M{"_id": oid}); err != nil {
			sc.AbortTransaction(context.Background())
			return err
		}

		sc.CommitTransaction(context.Background())
		return nil
	})
	if err != nil {
		return fmt.Errorf("fail 'repository.Recipe.DeleteOne': %w", err)
	}

	return nil
}
