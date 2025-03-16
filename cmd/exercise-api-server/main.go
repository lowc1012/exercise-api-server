package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lowc1012/exercise-api-server/internal/api"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "api-server",
	Short: "Code exercise",
	Long:  "HTTP API server exercise",
	RunE:  appRun,
}

func appRun(cmd *cobra.Command, args []string) error {
	log.Println("Starting API server")
	srv, err := api.StartAsync()
	if err != nil {
		log.Printf("API server failed to start, error: %v\n", err.Error())
		return err
	}

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown completed, error: %v\n", err.Error())
	}
	<-ctx.Done()
	log.Println("Server stopped successfully")
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
