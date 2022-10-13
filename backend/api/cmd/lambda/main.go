package main

import (
	"context"
	"doodemo-cook/api/database"
	"doodemo-cook/api/server"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("you must set your 'MONGODB_URI' environmental variable")
	}

	db, disconnect, err := database.New(context.Background(), uri)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer func() {
		if err := disconnect(context.Background()); err != nil {
			log.Printf("failed to disconnect db: %v", err)
		}
	}()

	r := server.NewRouter(db)

	lambda.Start(httpadapter.New(r).ProxyWithContext)
}
