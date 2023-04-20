package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ShadrackAdwera/go-subscriptions/api"
	db "github.com/ShadrackAdwera/go-subscriptions/db/sqlc"
	"github.com/ShadrackAdwera/go-subscriptions/token"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}
	symmetricKey := os.Getenv("SYMMETRIC_KEY")

	paseto, err := token.NewPasetoMaker(symmetricKey)

	if err != nil {
		log.Fatal(err)
	}

	url := os.Getenv("PG_TEST_URL")

	conn, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)

	srv := api.NewServer(store, paseto)

	err = srv.StartServer("0.0.0.0:5000")

	fmt.Println("Auth Service listening on PORT: 5000")

	if err != nil {
		log.Fatal(err)
	}
}
