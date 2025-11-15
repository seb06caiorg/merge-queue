package middleware

import (
	"net/http"
	"time"

	"merge-queue/internal/config"
	"merge-queue/pkg/utils"
)

// LoggingMiddleware logs HTTP requests.
type LoggingMiddleware struct {
	config *config.Config
	logger *utils.Logger
}

// NewLoggingMiddleware creates a new logging middleware instance.
func NewLoggingMiddleware(cfg *config.Config, logger *utils.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		config: cfg,
		logger: logger,
	}
}

// Handler returns the logging middleware handler.
func (lm *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !lm.config.Features.EnableLogging {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// Wrap the response writer to capture status code.
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		lm.logger.Info(
			"%s %s %d %v %s",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// DetailedLoggingMiddleware provides more detailed request logging.
type DetailedLoggingMiddleware struct {
	logger *utils.Logger
}

// NewDetailedLoggingMiddleware creates a detailed logging middleware.
func NewDetailedLoggingMiddleware(logger *utils.Logger) *DetailedLoggingMiddleware {
	return &DetailedLoggingMiddleware{logger: logger}
}

// Handler returns the detailed logging middleware handler.
func (dlm *DetailedLoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request details.
		dlm.logger.Debug(
			"Request started: %s %s from %s, User-Agent: %s",
			r.Method,
			r.URL.String(),
			r.RemoteAddr,
			r.Header.Get("User-Agent"),
		)

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Log response details.
		dlm.logger.Info(
			"Request completed: %s %s %d %v",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
		)

		if duration > 1*time.Second {
			dlm.logger.Warn(
				"Slow request detected: %s %s took %v",
				r.Method,
				r.URL.Path,
				duration,
			)
		}
	})
}
