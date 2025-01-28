package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Monkhai/shwipe-server.git/pkg/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer cancel()
	var wg sync.WaitGroup
	s, err := server.NewServer(ctx, &wg)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	http.HandleFunc("/ws", s.WebSocketHandler)

	go func() {
		log.Println("Starting server on port 8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	<-signalChan
	log.Println("Shutting down server...")

	cancel()
	log.Println("Waiting for all goroutines to finish...")
	wg.Wait()
	log.Println("Server shutdown complete")
}
