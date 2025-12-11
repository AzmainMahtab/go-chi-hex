package main

import (
	// ... imports for context, log, net/http, os, signal, syscall, time ...
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AzmainMahtab/docpad/api/http/handlers"
	routes "github.com/AzmainMahtab/docpad/api/http/router"
)

func main() {
	// ... Configuration (Port) setup ...
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	// Handler initiations
	healthHandler := handlers.NewHealthHandleer()

	// Router configuration
	deps := routes.RouterDependencies{
		HealthH: healthHandler,
	}

	router := routes.NewRouter(deps)

	// ... Server configuration (timeouts, handler) ...
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// ... Graceful shutdown logic ...
	go func() {
		log.Printf("🚀 DocPad API Server starting on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: Could not listen on %s: %v", addr, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("FATAL: Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting gracefully.")
}
