package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework/internal/kafka"
	"homework/internal/msgserv"
	"homework/internal/pb"
	"homework/internal/repository"
	"log"
	"os"
	"os/signal"
	"strings"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	conn, err := grpc.Dial(
		os.Getenv("GRPC_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("grpc: %v\n", err)
		return
	}

	msgClient := pb.NewMessageStorageClient(conn)
	userClient := pb.NewUserStorageClient(conn)

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	prod, err := kafka.NewProducer(brokers, "messages")
	if err != nil {
		log.Printf("kafka: %v\n", err)
		return
	}

	msgRepo := repository.NewMessageRepository(msgClient, prod)
	userRepo := repository.NewUserRepository(userClient)

	serv := msgserv.NewMsgServ(msgRepo, userRepo)
	err = serv.Serve(ctx, ":8080")
	if err != nil {
		log.Println(err)
	}
}
