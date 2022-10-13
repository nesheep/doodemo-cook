package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const name = "doodemo_cook"

func New(ctx context.Context) (*mongo.Database, func(ctx context.Context) error, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("you must set your 'MONGODB_URI' environmental variable")
	}

	opts := options.Client().ApplyURI(uri)
	cli, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	if err := cli.Ping(ctx, readpref.Primary()); err != nil {
		return nil, cli.Disconnect, err
	}

	db := cli.Database(name)
	return db, cli.Disconnect, nil
}
