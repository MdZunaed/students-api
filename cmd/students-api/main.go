package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MdZunaed/students-api/internal/config"
	"github.com/MdZunaed/students-api/internal/http/handlers/student"
	"github.com/MdZunaed/students-api/internal/storage/sqlite"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Database setup

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	// Setup router

	router := http.NewServeMux()

	router.HandleFunc("GET /api/students", student.Get())
	router.HandleFunc("POST /api/students", student.New(storage))

	// Setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	/// Graceful shutdown: If any issue with server,
	/// Server will shutdown after the current request. no instantly

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// err:= server.Shutdown(ctx)
	// if err != nil {
	// 	slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	// }
	// Shortcut:
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown successfully")
}
