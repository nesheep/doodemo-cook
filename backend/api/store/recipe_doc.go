package store

import (
	"doodemo-cook/api/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type insRecipe struct {
	Title string               `bson:"title"`
	URL   string               `bson:"url"`
	Tags  []primitive.ObjectID `bson:"tags"`
}

func newInsRecipe(r entity.Recipe) (insRecipe, error) {
	ir := insRecipe{
		Title: r.Title,
		URL:   r.URL,
		Tags:  []primitive.ObjectID{},
	}
	for _, v := range r.Tags {
		objId, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return insRecipe{}, err
		}
		ir.Tags = append(ir.Tags, objId)
	}
	return ir, nil
}
