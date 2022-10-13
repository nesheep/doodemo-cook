package store

import (
	"context"
	"doodemo-cook/api/entity"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const tagColl = "tags"

type Tag struct {
	db *mongo.Database
}

func NewTag(db *mongo.Database) *Tag {
	return &Tag{db: db}
}

func (s *Tag) Find(ctx context.Context) (entity.Tags, error) {
	coll := s.db.Collection(tagColl)

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return entity.Tags{}, err
	}

	data := []entity.Tag{}
	for cur.Next(ctx) {
		var tag entity.Tag
		cur.Decode(&tag)
		data = append(data, tag)
	}

	return entity.Tags{Data: data, Total: len(data)}, nil
}

func (s *Tag) FindOne(ctx context.Context, id string) (entity.Tag, error) {
	coll := s.db.Collection(tagColl)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.Tag{}, err
	}

	var tag entity.Tag
	if err := coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&tag); err != nil {
		return entity.Tag{}, err
	}

	return tag, nil
}

func (s *Tag) InsertOne(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	coll := s.db.Collection(tagColl)

	result, err := coll.InsertOne(ctx, tag)
	if err != nil {
		return entity.Tag{}, err
	}

	id := result.InsertedID
	var newTag entity.Tag
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&newTag); err != nil {
		return entity.Tag{}, err
	}

	return newTag, nil
}

func (s *Tag) UpdateOne(ctx context.Context, id string, tag entity.Tag) (entity.Tag, error) {
	coll := s.db.Collection(tagColl)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.Tag{}, err
	}

	result, err := coll.ReplaceOne(ctx, bson.M{"_id": objId}, tag)
	if err != nil {
		return entity.Tag{}, err
	}

	if result.MatchedCount < 1 {
		return entity.Tag{}, fmt.Errorf("not found %s", id)
	}

	var newTag entity.Tag
	if err := coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&newTag); err != nil {
		return entity.Tag{}, err
	}

	return newTag, nil
}

func (s *Tag) DeleteOne(ctx context.Context, id string) (string, error) {
	coll := s.db.Collection(tagColl)

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
