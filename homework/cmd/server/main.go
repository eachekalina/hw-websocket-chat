package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/msgserv"
	"homework/internal/repository"
	"log"
	"os"
	"os/signal"
)

func main() {
	connString, found := os.LookupEnv("DB_URL")
	if !found {
		connString = "postgres://postgres:mysecretpassword@localhost:5432/postgres"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	msgRepo := repository.NewMessageRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	serv := msgserv.NewMsgServ(msgRepo, userRepo)
	err = serv.Serve(ctx, ":8080")
	if err != nil {
		log.Println(err)
	}
}
