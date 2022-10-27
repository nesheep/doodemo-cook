package main

import (
	"context"
	"doodemo-cook/api/database"
	"doodemo-cook/api/server"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

var handlerAdapter *httpadapter.HandlerAdapter

func init() {
	c, err := database.New(context.Background())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	r := server.NewRouter(c)
	handlerAdapter = httpadapter.New(r)
}

func main() {
	lambda.Start(handlerAdapter.ProxyWithContext)
}
