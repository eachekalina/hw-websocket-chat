package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	sigQuit := make(chan os.Signal, 2)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			fmt.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	eg.Go(func() error {
		count := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Println("graceful shutdown process")
				return fmt.Errorf("done")
			default:
				// здесь может быть что угодно, например обработка событий с кафки
				count++
			}
			fmt.Println(count)
			time.Sleep(time.Second)
		}
	})

	_ = eg.Wait()
	fmt.Println("graceful shutdown service")
	//_ = server.Shutdown(context.Background())
}
