package repository

import (
	"context"
	"doodemo-cook/api/entity"
	"errors"

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

func (r *Tag) Count(ctx context.Context) (int, error) {
	coll := r.db.Collection(tagColl)

	cnt, err := coll.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}

	return int(cnt), nil
}

func (r *Tag) Find(ctx context.Context) (entity.Tags, error) {
	coll := r.db.Collection(tagColl)

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	tags := entity.Tags{}
	for cur.Next(ctx) {
		var b bTag
		cur.Decode(&b)
		tags = append(tags, b.toTag())
	}

	return tags, nil
}

func (r *Tag) findByName(ctx context.Context, name string) (entity.Tag, error) {
	coll := r.db.Collection(tagColl)

	var b bTag
	if err := coll.FindOne(ctx, bson.M{"name": name}).Decode(&b); err != nil {
		return entity.Tag{}, err
	}

	return b.toTag(), nil
}

func (r *Tag) insertOne(ctx context.Context, tag entity.Tag) (string, error) {
	coll := r.db.Collection(tagColl)

	b := bInputTagFromTag(tag)
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

func (r *Tag) insertIfNotExists(ctx context.Context, tag entity.Tag) (string, error) {
	existTag, err := r.findByName(ctx, tag.Name)
	if err == mongo.ErrNoDocuments {
		id, err := r.insertOne(ctx, tag)
		if err != nil {
			return "", err
		}
		return id, nil
	}
	if err != nil {
		return "", err
	}

	return existTag.ID, nil
}
