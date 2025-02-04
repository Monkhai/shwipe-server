package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Monkhai/shwipe-server.git/pkg/server"
)

// func printUserGoroutines() {
// 	buf := make([]byte, 1<<16)
// 	runtime.Stack(buf, true)
// 	stacks := string(buf)

// 	goroutines := strings.Split(stacks, "\n\n")

// 	log.Printf("=== User Goroutines ===\n")
// 	for _, g := range goroutines {
// 		if strings.Contains(g, "runtime.") ||
// 			strings.Contains(g, "system") ||
// 			strings.Contains(g, "GC") ||
// 			strings.Contains(g, "finalizer") {
// 			continue
// 		}

// 		if strings.TrimSpace(g) == "" {
// 			continue
// 		}

// 		if strings.Contains(g, "pkg/session") ||
// 			strings.Contains(g, "pkg/server") {
// 			log.Printf("%s\n", g)
// 		}
// 	}
// 	log.Printf("=====================\n")
// }

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

	// go func() {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		case <-time.After(5 * time.Second):
	// 			log.Printf("Current number of goroutines: %d", runtime.NumGoroutine())
	// 			printUserGoroutines()
	// 		}
	// 	}
	// }()

	<-signalChan

	log.Println("Shutting down server...")
	cancel()
	log.Println("Waiting for all goroutines to finish...")

	log.Printf("WaitGroup state before wait")
	// printUserGoroutines()

	waitChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-waitChan:
		log.Println("WaitGroup completed normally")
	case <-time.After(10 * time.Second):
		log.Println("WARNING: WaitGroup wait timed out!")
		// printUserGoroutines()
	}

	log.Println("Server shutdown complete")
}
