package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/m4rk0G/keydra/internal/unsealer"
	"github.com/m4rk0G/keydra/pkg/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "/etc/keydra/config.yaml", "Path to configuration file")
	flag.Parse()

	log.Println("Starting Keydra - HashiCorp Vault Auto-Unsealer")

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v, shutting down gracefully...", sig)
		cancel()
	}()

	var wg sync.WaitGroup
	for _, _ = range cfg.Vault.Nodes {
		wg.Add(1)

		// Create unsealer instance
		unsealer, err := unsealer.New(cfg)
		if err != nil {
			log.Fatalf("Failed to create unsealer: %v", err)
		}

		// Start the unsealer
		log.Println("Starting Vault auto-unsealer...")
		go unsealer.Start(ctx, &wg)
	}

	wg.Wait()

	log.Println("Keydra shutdown complete")
}
