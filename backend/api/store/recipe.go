package store

import (
	"context"
	"doodemo-cook/api/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const name = "recipes"

type Recipe struct {
	db *mongo.Database
}

func NewRecipe(db *mongo.Database) *Recipe {
	return &Recipe{db: db}
}

type recipeBSON struct {
	ID    string `bson:"_id"`
	Title string `bson:"title"`
	URL   string `bson:"URL"`
}

type insertRecipeBSON struct {
	Title string `bson:"title"`
	URL   string `bson:"URL"`
}

func (s *Recipe) Find(ctx context.Context) ([]entity.Recipe, error) {
	coll := s.db.Collection(name)

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	recipes := []entity.Recipe{}

	for cur.Next(ctx) {
		var data recipeBSON
		cur.Decode(&data)
		recipes = append(recipes, entity.Recipe{
			ID:    data.ID,
			Title: data.Title,
			URL:   data.URL,
		})
	}

	return recipes, nil
}

func (s *Recipe) InsertOne(ctx context.Context, recipe entity.Recipe) (entity.Recipe, error) {
	coll := s.db.Collection(name)

	insertDoc := insertRecipeBSON{
		Title: recipe.Title,
		URL:   recipe.URL,
	}

	result, err := coll.InsertOne(ctx, insertDoc)
	if err != nil {
		return entity.Recipe{}, err
	}

	id := result.InsertedID
	var insertedDoc recipeBSON
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&insertedDoc); err != nil {
		return entity.Recipe{}, err
	}

	return entity.Recipe{
		ID:    insertedDoc.ID,
		Title: insertedDoc.Title,
		URL:   insertedDoc.URL,
	}, nil
}
