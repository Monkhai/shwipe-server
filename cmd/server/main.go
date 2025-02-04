package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Monkhai/shwipe-server.git/pkg/server"
)

func printUserGoroutines() {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	stacks := string(buf)

	goroutines := strings.Split(stacks, "\n\n")

	log.Printf("=== User Goroutines ===\n")
	for _, g := range goroutines {
		if strings.Contains(g, "runtime.") ||
			strings.Contains(g, "system") ||
			strings.Contains(g, "GC") ||
			strings.Contains(g, "finalizer") {
			continue
		}

		if strings.TrimSpace(g) == "" {
			continue
		}

		if strings.Contains(g, "pkg/session") ||
			strings.Contains(g, "pkg/db") ||
			strings.Contains(g, "pkg/server") {
			log.Printf("%s\n", g)
			log.Println("--------------------------------")
		}
	}
	log.Printf("=====================\n")
}

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
	http.HandleFunc("/get-sessions", s.GetSessions)
	http.HandleFunc("/get-user", s.GetUser)
	http.HandleFunc("/get-users", s.GetUsers)
	go func() {
		log.Println("Starting server on port 8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	<-signalChan
	cancel()

	waitChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitChan)
	}()

	go func() {
		if err := s.Shutdown(); err != nil {
			log.Printf("HTTP server Shutdown error: %v", err)
		}
	}()

	select {
	case <-waitChan:
		log.Println("WaitGroup completed normally")
	case <-time.After(10 * time.Second):
		log.Println("WARNING: WaitGroup wait timed out!")
		printUserGoroutines()
	}

	log.Println("Server shutdown complete")
}
