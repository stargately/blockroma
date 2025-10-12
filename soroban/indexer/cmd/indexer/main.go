package main

import (
	"context"
	"flag"
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
	// Parse CLI flags
	startLedger := flag.Uint("start-ledger", 0, "Start ledger for backfill mode (0 = live polling)")
	endLedger := flag.Uint("end-ledger", 0, "End ledger for backfill mode (0 = current ledger)")
	batchSize := flag.Uint("batch-size", 100, "Batch size for backfill mode")
	rateLimit := flag.Uint("rate-limit", 10, "Max requests per second for backfill mode")
	flag.Parse()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Set log level from environment variable
	logLevel := getEnv("LOG_LEVEL", "info")
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Get configuration from environment (minimal env vars)
	rpcURL := getEnv("STELLAR_RPC_URL", "http://stellar-rpc:8000")
	postgresURL := getEnv("POSTGRES_DSN", "")

	if postgresURL == "" {
		logger.Fatal("POSTGRES_DSN environment variable is required")
	}

	// Determine mode
	isBackfillMode := *startLedger > 0
	if isBackfillMode {
		logger.WithFields(logrus.Fields{
			"startLedger": *startLedger,
			"endLedger":   *endLedger,
			"batchSize":   *batchSize,
			"rateLimit":   *rateLimit,
		}).Info("Starting in BACKFILL mode")
	} else {
		logger.Info("Starting in LIVE POLLING mode")
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

	// Start appropriate mode
	errCh := make(chan error, 1)
	if isBackfillMode {
		// Backfill mode: process historical ledgers
		go func() {
			config := poller.BackfillConfig{
				StartLedger: uint32(*startLedger),
				EndLedger:   uint32(*endLedger),
				BatchSize:   uint32(*batchSize),
				RateLimit:   uint32(*rateLimit),
			}
			if err := p.Backfill(ctx, config); err != nil {
				errCh <- err
			} else {
				logger.Info("Backfill completed successfully")
				cancel() // Exit after successful backfill
			}
		}()
	} else {
		// Live polling mode: continuous processing
		go func() {
			if err := p.Start(ctx); err != nil {
				errCh <- err
			}
		}()
	}

	// Wait for shutdown signal or error
	select {
	case <-sigCh:
		logger.Info("Received shutdown signal")
		cancel()
		time.Sleep(2 * time.Second) // Grace period
	case err := <-errCh:
		if err != nil {
			logger.WithError(err).Error("Indexer error")
		}
		cancel()
	case <-ctx.Done():
		// Context cancelled (e.g., backfill completed)
		logger.Info("Context cancelled")
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
