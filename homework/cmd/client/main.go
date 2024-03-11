package main

import (
	"context"
	"homework/internal/msgclient"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := msgclient.NewMsgClient()
	err := client.Run(ctx, "localhost:8080")
	if err != nil {
		log.Println(err)
	}
}
