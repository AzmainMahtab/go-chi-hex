// @title DocPad Hospital Management API
// @version 1.0
// @description The core API for DocPad, managing patient records, scheduling, and prescriptions for the Bangladeshi demographic.
// @contact.name DocPad Support
// @contact.email support@docpad.bd
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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

	"github.com/AzmainMahtab/go-chi-hex/api/http/handlers"
	routes "github.com/AzmainMahtab/go-chi-hex/api/http/router"
	"github.com/AzmainMahtab/go-chi-hex/internal/config"
	"github.com/AzmainMahtab/go-chi-hex/internal/infrastructure/postgres"
	"github.com/AzmainMahtab/go-chi-hex/internal/services/users"
)

func main() {
	log.Println("Starting DocPad service assembly...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("FATAL: Failed to load configuration: %v", err)
	}

	dbConfig := postgres.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		PoolSize: cfg.DB.PoolSize,
	}

	// connect to the Database
	db, err := postgres.ConnectDB(dbConfig)
	if err != nil {
		log.Fatalf("FATAL: Database connection failed: %v", err)
	}
	defer db.Close() // Ensure the connection is closed on exit

	// REPOSITORY SETUP
	userRepo := postgres.NewUserRepo(db)

	// SERVICE SETUP
	userService := users.NewUserService(userRepo)

	// HANDLER AND ROUTER SETUP
	healthHandler := handlers.NewHealthHandleer()
	userHandler := handlers.NewUserHandler(userService)

	deps := routes.RouterDependencies{
		HealthH: healthHandler,
		UserH:   userHandler,
	}
	router := routes.NewRouter(deps)

	// SERVER SETUP
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server in a non-blocking goroutine
	go func() {
		log.Printf("🚀 DocPad API Server starting on http://localhost:%s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: Could not listen on %s: %v", cfg.Server.Port, err)
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
