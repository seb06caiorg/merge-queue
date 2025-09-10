package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"merge-queue/internal/config"
	"merge-queue/internal/handlers"
	"merge-queue/internal/middleware"
	"merge-queue/internal/services"
	"merge-queue/pkg/utils"
)

func main() {
	// Load configuration.
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger.
	logLevel := utils.InfoLevel
	if cfg.App.Debug {
		logLevel = utils.DebugLevel
	}
	logger := utils.NewLogger(logLevel)

	logger.Info("Starting %s v%s", cfg.App.Name, cfg.App.Version)
	logger.Info("Environment: %s", cfg.App.Environment)

	// Initialize services.
	taskService := services.NewTaskService(cfg.Features.MaxTasksPerUser)

	// Initialize handlers.
	taskHandler := handlers.NewTaskHandler(taskService, logger)
	healthHandler := handlers.NewHealthHandler(cfg, logger)
	staticHandler := handlers.NewStaticHandler(cfg, logger)

	// Initialize middleware.
	corsMiddleware := middleware.NewCORSMiddleware(cfg)
	loggingMiddleware := middleware.NewLoggingMiddleware(cfg, logger)
	authMiddleware := middleware.NewAuthMiddleware(logger)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(cfg, logger)

	// Setup router.
	router := setupRouter(
		taskHandler,
		healthHandler,
		staticHandler,
		corsMiddleware,
		loggingMiddleware,
		authMiddleware,
		rateLimitMiddleware,
	)

	// Create HTTP server.
	server := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine.
	go func() {
		logger.Info("🚀 Server starting on http://localhost%s", cfg.Server.Port)
		logger.Info("📋 Sample tasks loaded and ready for your hackathon!")
		logger.Info("🌐 Web interface: http://localhost%s", cfg.Server.Port)
		logger.Info("📖 API docs: http://localhost%s/api/v1/health", cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server.
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	// Cleanup middleware.
	rateLimitMiddleware.Stop()

	logger.Info("Server gracefully stopped")
}

// setupRouter configures and returns the HTTP router.
func setupRouter(
	taskHandler *handlers.TaskHandler,
	healthHandler *handlers.HealthHandler,
	staticHandler *handlers.StaticHandler,
	corsMiddleware *middleware.CORSMiddleware,
	loggingMiddleware *middleware.LoggingMiddleware,
	authMiddleware *middleware.AuthMiddleware,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
) *mux.Router {
	router := mux.NewRouter()

	// Apply global middleware.
	router.Use(corsMiddleware.Handler)
	router.Use(loggingMiddleware.Handler)
	router.Use(rateLimitMiddleware.Handler)

	// API routes.
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health endpoints (no auth required).
	api.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")
	api.HandleFunc("/ready", healthHandler.ReadinessCheck).Methods("GET")
	api.HandleFunc("/live", healthHandler.LivenessCheck).Methods("GET")

	// Task endpoints (with optional auth).
	api.Use(authMiddleware.Handler) // Optional auth for all API routes.

	// Task CRUD operations.
	api.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET")
	api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.GetTask).Methods("GET")
	api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.UpdateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.DeleteTask).Methods("DELETE")

	// Additional task operations.
	api.HandleFunc("/tasks/search", taskHandler.SearchTasks).Methods("POST")
	api.HandleFunc("/tasks/stats", taskHandler.GetTaskStats).Methods("GET")

	// Static content.
	router.HandleFunc("/", staticHandler.ServeHome).Methods("GET")

	// Handle 404s with a custom response.
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := utils.NewResponseHelper()
		response.SendError(w, http.StatusNotFound, fmt.Sprintf("Endpoint not found: %s %s", r.Method, r.URL.Path))
	})

	return router
}
