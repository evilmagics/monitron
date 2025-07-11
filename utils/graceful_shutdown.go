package utils

import (
	"context"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// WaitForShutdownSignal waits for OS signals to initiate graceful shutdown.
func WaitForShutdownSignal(cleanup func(context.Context)) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // Block until a signal is received.

	log.Info().Msg("Received shutdown signal. Initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 10-second timeout for cleanup
	defer cancel()

	cleanup(ctx)

	log.Info().Msg("Graceful shutdown complete. Exiting.")
	os.Exit(0)
}

