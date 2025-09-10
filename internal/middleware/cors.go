package middleware

import (
	"fmt"
	"net/http"

	"merge-queue/internal/config"
)

// CORSMiddleware handles Cross-Origin Resource Sharing.
type CORSMiddleware struct {
	config *config.Config
}

// NewCORSMiddleware creates a new CORS middleware instance.
func NewCORSMiddleware(cfg *config.Config) *CORSMiddleware {
	return &CORSMiddleware{config: cfg}
}

// Handler returns the CORS middleware handler.
func (cm *CORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cm.config.Features.EnableCORS {
			next.ServeHTTP(w, r)
			return
		}

		// Set CORS headers.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours.

		// Handle preflight requests.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ConfigurableCORSMiddleware allows more fine-grained CORS control.
type ConfigurableCORSMiddleware struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

// NewConfigurableCORSMiddleware creates a configurable CORS middleware.
func NewConfigurableCORSMiddleware(origins, methods, headers []string, maxAge int) *ConfigurableCORSMiddleware {
	return &ConfigurableCORSMiddleware{
		AllowedOrigins: origins,
		AllowedMethods: methods,
		AllowedHeaders: headers,
		MaxAge:         maxAge,
	}
}

// Handler returns the configurable CORS middleware handler.
func (ccm *ConfigurableCORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if origin is allowed.
		allowed := false
		for _, allowedOrigin := range ccm.AllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Set other CORS headers.
		if len(ccm.AllowedMethods) > 0 {
			methods := ""
			for i, method := range ccm.AllowedMethods {
				if i > 0 {
					methods += ", "
				}
				methods += method
			}
			w.Header().Set("Access-Control-Allow-Methods", methods)
		}

		if len(ccm.AllowedHeaders) > 0 {
			headers := ""
			for i, header := range ccm.AllowedHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			w.Header().Set("Access-Control-Allow-Headers", headers)
		}

		if ccm.MaxAge > 0 {
			w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", ccm.MaxAge))
		}

		// Handle preflight requests.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
