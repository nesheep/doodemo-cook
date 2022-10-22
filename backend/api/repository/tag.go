package repository

import (
	"context"
	"doodemo-cook/api/entity"

	"go.mongodb.org/mongo-driver/mongo"
)

type Tag struct {
	db *mongo.Database
}

func NewTag(db *mongo.Database) *Tag {
	return &Tag{db: db}
}

func (r *Tag) Count(ctx context.Context) (int, error) {
	return 0, nil
}

func (r *Tag) Find(ctx context.Context) (entity.Tags, error) {
	tags := entity.Tags{}
	return tags, nil
}
