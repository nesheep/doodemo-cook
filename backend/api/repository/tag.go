package repository

import (
	"context"
	"doodemo-cook/api/entity"

	"go.mongodb.org/mongo-driver/bson"
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
		return entity.Tags{}, err
	}

	tags := entity.Tags{}
	for cur.Next(ctx) {
		var t bTag
		cur.Decode(&t)
		tags = append(tags, t.toTag())
	}

	return tags, nil
}
