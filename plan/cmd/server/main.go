package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"plan/internal/db"
	"plan/internal/server"
)

func main() {
	// Define flags
	port := flag.Int("port", 8080, "Server port")
	host := flag.String("host", "localhost", "Server host")
	dev := flag.Bool("dev", false, "Development mode (proxy to Vite dev server)")
	webDir := flag.String("web-dir", "../web/dist", "Path to web dist directory")

	flag.Parse()

	// Initialize database
	if err := db.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Format address
	addr := fmt.Sprintf("%s:%d", *host, *port)

	// Print startup info
	fmt.Printf("Starting server at http://%s\n", addr)
	fmt.Printf("WebSocket endpoint: ws://%s/api/ws\n", addr)

	if *dev {
		fmt.Println("Development mode: proxying to Vite dev server at http://localhost:5173")
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start(addr, *dev, *webDir)
	}()

	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		fmt.Printf("\nShutdown signal received (%v), closing connections...\n", sig)
		// Database cleanup happens via defer
	case err := <-errChan:
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
