package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const name = "doodemo_cook"

func New(ctx context.Context, uri string) (*mongo.Database, func(ctx context.Context) error, error) {
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
