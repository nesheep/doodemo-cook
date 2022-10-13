package main

import (
	"context"
	"doodemo-cook/api/database"
	"doodemo-cook/api/server"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context) error {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	db, disconnect, err := database.New(ctx)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer func() {
		if err := disconnect(context.Background()); err != nil {
			log.Printf("failed to disconnect db: %v", err)
		}
	}()

	port := 8888
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	r := server.NewRouter(db)
	s := server.NewServer(r, l)

	return s.Run(ctx)
}
