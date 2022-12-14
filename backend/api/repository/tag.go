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

type Tag struct {
	c *mongo.Client
}

func NewTag(c *mongo.Client) *Tag {
	return &Tag{c: c}
}

func (r *Tag) coll() *mongo.Collection {
	return r.c.Database(dbName).Collection(tagColl)
}

func (r *Tag) Count(ctx context.Context) (int, error) {
	coll := r.coll()

	cnt, err := coll.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, fmt.Errorf("fail 'repository.Tag.Count': %w", err)
	}

	return int(cnt), nil
}

func (r *Tag) Find(ctx context.Context) (entity.Tags, error) {
	coll := r.coll()

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("fail 'repository.Tag.Find': %w", err)
	}

	tags := entity.Tags{}
	for cur.Next(ctx) {
		var b bTag
		cur.Decode(&b)
		tags = append(tags, b.toTag())
	}

	return tags, nil
}

func (r *Tag) findOne(ctx context.Context, id string) (entity.Tag, error) {
	coll := r.coll()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.Tag{}, fmt.Errorf("fail 'repository.Tag.FindOne': %w", err)
	}

	var b bTag
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&b); err != nil {
		return entity.Tag{}, fmt.Errorf("fail 'repository.Tag.FindOne': %w", err)
	}

	return b.toTag(), nil
}

func (r *Tag) findByName(ctx context.Context, name string) (entity.Tag, error) {
	coll := r.coll()

	var b bTag
	if err := coll.FindOne(ctx, bson.M{"name": name}).Decode(&b); err != nil {
		return entity.Tag{}, fmt.Errorf("fail 'repository.Tag.findByName': %w", err)
	}

	return b.toTag(), nil
}

func (r *Tag) insertOne(ctx context.Context, tag entity.Tag) (string, error) {
	coll := r.coll()

	b := bInputTagFromTag(tag)
	b.RecipeNum = 1
	result, err := coll.InsertOne(ctx, b)
	if err != nil {
		return "", fmt.Errorf("fail 'repository.Tag.insertOne': %w", err)
	}

	id := result.InsertedID
	oid, ok := id.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("fail 'repository.Tag.insertOne': %w", err)
	}

	return oid.Hex(), nil
}

func (r *Tag) incrementRecipeNum(ctx context.Context, id string, amount int) error {
	coll := r.coll()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("fail 'repository.Tag.incrementRecipeNum': %w", err)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$inc": bson.M{"recipe_num": amount}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("fail 'repository.Tag.incrementRecipeNum': %w", err)
	}

	if result.MatchedCount < 1 {
		return fmt.Errorf("fail 'repository.Tag.incrementRecipeNum': %w", err)
	}

	return nil
}

func (r *Tag) deleteOne(ctx context.Context, id string) error {
	coll := r.coll()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("fail 'repository.Tag.deleteOne': %w", err)
	}

	result, err := coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("fail 'repository.Tag.deleteOne': %w", err)
	}

	if result.DeletedCount < 1 {
		return fmt.Errorf("fail 'repository.Tag.deleteOne': %w", err)
	}

	return nil
}

func (r *Tag) insertOrIncrement(ctx context.Context, tag entity.Tag) (string, error) {
	existTag, err := r.findByName(ctx, tag.Name)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			id, err := r.insertOne(ctx, tag)
			if err != nil {
				return "", fmt.Errorf("fail 'repository.Tag.insertOrIncremet': %w", err)
			}
			return id, nil
		}
		return "", fmt.Errorf("fail 'repository.Tag.insertOrIncremet': %w", err)
	}

	if err := r.incrementRecipeNum(ctx, existTag.ID, 1); err != nil {
		return "", fmt.Errorf("fail 'repository.Tag.insertOrIncremet': %w", err)
	}

	return existTag.ID, nil
}

func (r *Tag) deleteOrDecrement(ctx context.Context, id string) error {
	tag, err := r.findOne(ctx, id)
	if err != nil {
		return fmt.Errorf("fail 'repository.Tag.deleteOrDecrement': %w", err)
	}

	if tag.RecipeNum <= 1 {
		if err := r.deleteOne(ctx, id); err != nil {
			return fmt.Errorf("fail 'repository.Tag.deleteOrDecrement': %w", err)
		}
		return nil
	}

	if err := r.incrementRecipeNum(ctx, id, -1); err != nil {
		return fmt.Errorf("fail 'repository.Tag.deleteOrDecrement': %w", err)
	}

	return nil
}
