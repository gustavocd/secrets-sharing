package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gustavocd/secrets-sharing/filestore"
	"github.com/gustavocd/secrets-sharing/handlers"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	fPath := os.Getenv("DATA_FILE_PATH")
	if fPath == "" {
		return fmt.Errorf("DATA_FILE_PATH variable not set")
	}

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	err := filestore.Init(fPath)
	if err != nil {
		return fmt.Errorf("could not initialize file store: %w", err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apiMux := handlers.SetupHandlers(handlers.APIMuxConfig{
		Shutdown: shutdown,
	})

	api := http.Server{
		Addr:    addr,
		Handler: apiMux,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Println("starting server on ", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Println("shutdown started with signal", sig)

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
