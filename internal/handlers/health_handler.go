package handlers

import (
	"net/http"
	"time"

	"merge-queue/internal/config"
	"merge-queue/internal/models"
	"merge-queue/pkg/utils"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	config    *config.Config
	response  *utils.ResponseHelper
	logger    *utils.Logger
	startTime time.Time
}

// NewHealthHandler creates a new HealthHandler instance.
func NewHealthHandler(cfg *config.Config, logger *utils.Logger) *HealthHandler {
	return &HealthHandler{
		config:    cfg,
		response:  utils.NewResponseHelper(),
		logger:    logger,
		startTime: time.Now(),
	}
}

// HealthCheck handles GET /health requests.
func (hh *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(hh.startTime)

	response := models.HealthResponse{
		Status:    "healthy",
		Version:   hh.config.App.Version,
		Timestamp: time.Now(),
		Uptime:    utils.NewTimeUtils().FormatDuration(uptime),
	}

	hh.response.SendSuccess(w, response)
}

// ReadinessCheck handles GET /ready requests.
func (hh *HealthHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	// In a real application, you'd check database connectivity,
	// external service availability, etc.

	checks := map[string]string{
		"database":     "ok", // Placeholder.
		"external_api": "ok", // Placeholder.
		"memory":       "ok", // Could check memory usage.
		"disk":         "ok", // Could check disk space.
	}

	allHealthy := true
	for _, status := range checks {
		if status != "ok" {
			allHealthy = false
			break
		}
	}

	response := map[string]interface{}{
		"status": func() string {
			if allHealthy {
				return "ready"
			}
			return "not_ready"
		}(),
		"checks":    checks,
		"timestamp": time.Now(),
	}

	statusCode := http.StatusOK
	if !allHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	w.WriteHeader(statusCode)
	hh.response.SendJSON(w, statusCode, response)
}

// LivenessCheck handles GET /live requests.
func (hh *HealthHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	// Simple liveness check - if we can respond, we're alive.
	response := map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now(),
		"uptime":    utils.NewTimeUtils().FormatDuration(time.Since(hh.startTime)),
	}

	hh.response.SendSuccess(w, response)
}
