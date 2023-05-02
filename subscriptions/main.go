package main

import (
	"database/sql"
	"log"
	"os"

	api "github.com/ShadrackAdwera/go-subscriptions/subscriptions/api"
	db "github.com/ShadrackAdwera/go-subscriptions/subscriptions/db/sqlc"
	"github.com/ShadrackAdwera/go-subscriptions/subscriptions/workers"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"

	zerolog "github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	url := os.Getenv("PG_SUBSCRIPTIONS_URL")
	redisAddr := os.Getenv("REDIS_ADDRESS")

	conn, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)

	redisConf := &asynq.RedisClientOpt{
		Addr: redisAddr,
	}

	srv := api.NewServer(store)

	go ListenToQueue(redisConf, store)
	err = srv.StartServer("0.0.0.0:5001")

	if err != nil {
		panic(err)
	}

}

func ListenToQueue(opts *asynq.RedisClientOpt, store db.TxStore) {
	processor := workers.NewProcessor(opts, store)

	err := processor.Start()

	if err != nil {
		zerolog.Fatal().Err(err).Msg("error starting the redis task processor")
		return
	}
	zerolog.Info().Msg("redis task processor started . . . ")
}
