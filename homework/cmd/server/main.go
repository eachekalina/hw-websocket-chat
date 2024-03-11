package main

import (
	"context"
	"homework/internal/msgserv"
	"homework/internal/repository"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	repo, err := repository.NewRepo(ctx, "postgres://postgres:mysecretpassword@localhost:5432/postgres")
	if err != nil {
		log.Println(err)
		return
	}
	defer repo.Close()

	serv := msgserv.NewMsgServ(repo)
	err = serv.Serve(ctx, ":8080")
	if err != nil {
		log.Println(err)
	}
}
