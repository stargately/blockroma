package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/blockroma/soroban-indexer/pkg/client"
	"github.com/blockroma/soroban-indexer/pkg/db"
	"github.com/blockroma/soroban-indexer/pkg/poller"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// Get configuration from environment (minimal env vars)
	rpcURL := getEnv("STELLAR_RPC_URL", "http://stellar-rpc:8000")
	postgresURL := getEnv("POSTGRES_DSN", "")

	if postgresURL == "" {
		logger.Fatal("POSTGRES_DSN environment variable is required")
	}

	logger.WithFields(logrus.Fields{
		"rpcURL": rpcURL,
	}).Info("Starting Stellar RPC Indexer")

	// Connect to database
	database, err := db.Connect(postgresURL)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer database.Close()

	// Create RPC client
	rpcClient := client.NewClient(rpcURL)

	// Check RPC connectivity
	ctx := context.Background()
	if err := rpcClient.Health(ctx); err != nil {
		logger.WithError(err).Fatal("RPC health check failed")
	}

	network, err := rpcClient.GetNetwork(ctx)
	if err != nil {
		logger.WithError(err).Warn("Failed to get network info")
	} else {
		logger.WithFields(logrus.Fields{
			"passphrase": network.Passphrase,
			"protocol":   network.ProtocolVersion,
		}).Info("Connected to Stellar network")
	}

	// Create poller
	p := poller.New(rpcClient, database.DB, logger)

	// Start health/metrics HTTP server
	go startHTTPServer(p, logger)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start poller in goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := p.Start(ctx); err != nil {
			errCh <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-sigCh:
		logger.Info("Received shutdown signal")
		cancel()
		time.Sleep(2 * time.Second) // Grace period
	case err := <-errCh:
		logger.WithError(err).Error("Poller error")
		cancel()
	}

	logger.Info("Indexer stopped")
}

// startHTTPServer starts health and metrics HTTP server
func startHTTPServer(p *poller.Poller, logger *logrus.Logger) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		stats, err := p.GetStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"lastLedger": %v,
			"totalEvents": %v,
			"totalTransactions": %v
		}`, stats["lastLedger"], stats["totalEvents"], stats["totalTransactions"])
	})

	logger.Info("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.WithError(err).Error("HTTP server error")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
