package repository

import (
	"doodemo-cook/api/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type bInputTag struct {
	Name      string `bson:"name"`
	RecipeNum int    `bson:"recipe_num"`
}

func bInputTagFromTag(tag entity.Tag) bInputTag {
	return bInputTag{Name: tag.Name}
}

type bTag struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	RecipeNum int                `bson:"recipe_num"`
}

func (t bTag) toTag() entity.Tag {
	return entity.Tag{
		ID:        t.ID.Hex(),
		Name:      t.Name,
		RecipeNum: t.RecipeNum,
	}
}

type bTags []bTag

func (t bTags) toTags() entity.Tags {
	tags := make(entity.Tags, 0, len(t))
	for _, v := range t {
		tags = append(tags, v.toTag())
	}

	return tags
}
