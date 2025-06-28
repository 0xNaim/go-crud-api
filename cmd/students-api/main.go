package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xNaim/students-api/internal/config"
	"github.com/0xNaim/students-api/internal/http/handlers/student"
	"github.com/0xNaim/students-api/internal/storage/sqlite"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Initialize storage
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage initialized successfully", slog.String("storage_type", "sqlite"), slog.String("env", "development"), slog.String("storage_path", cfg.StoragePath))

	// Route setup
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetByID(storage))
	router.HandleFunc("GET /api/students", student.GetAll(storage))

	// Server setup
	server := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		fmt.Println("Server started successfully on", cfg.HTTPServer.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	fmt.Println("\nShutting down server gracefully...")

	// Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %v\n", err)
	} else {
		fmt.Println("Server stopped gracefully.")
	}
}
