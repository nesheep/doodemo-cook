package repository

import (
	"doodemo-cook/api/entity"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type bInputRecipe struct {
	Title string               `bson:"title"`
	URL   string               `bson:"url"`
	Tags  []primitive.ObjectID `bson:"tags"`
}

func bInputRecipeFromRecipe(recipe entity.Recipe) (bInputRecipe, error) {
	tags := make([]primitive.ObjectID, 0, len(recipe.Tags))
	for _, v := range recipe.Tags {
		oid, err := primitive.ObjectIDFromHex(v.ID)
		if err != nil {
			return bInputRecipe{}, fmt.Errorf("fail 'repository.bInputRecipeFromRecipe': %w", err)
		}
		tags = append(tags, oid)
	}

	return bInputRecipe{
		Title: recipe.Title,
		URL:   recipe.URL,
		Tags:  tags,
	}, nil
}

type bRecipe struct {
	ID    primitive.ObjectID `bson:"_id"`
	Title string             `bson:"title"`
	URL   string             `bson:"url"`
	Tags  bTags              `bson:"tags"`
}

func (r bRecipe) toRecipe() entity.Recipe {
	return entity.Recipe{
		ID:    r.ID.Hex(),
		Title: r.Title,
		URL:   r.URL,
		Tags:  r.Tags.toTags(),
	}
}
