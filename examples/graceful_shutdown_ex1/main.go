package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(&log.TextFormatter{})

	sigQuit := make(chan os.Signal, 2)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("captured signal: %v", s)
		}
	})

	server := &http.Server{}

	go func() {
		_ = server.ListenAndServe() // check error!
	}()

	if err := eg.Wait(); err != nil {
		logger.Infof("gracefully shutting down the server: %v", err)
	}

	_ = server.Shutdown(context.Background()) // check error!

	// для grpc сервера тут был бы grpcServer.GracefulStop()
}
