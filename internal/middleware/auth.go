package middleware

import (
	"context"
	"net/http"
	"strings"

	"merge-queue/pkg/utils"
)

// AuthMiddleware handles authentication (placeholder for future implementation).
type AuthMiddleware struct {
	logger *utils.Logger
}

// NewAuthMiddleware creates a new auth middleware instance.
func NewAuthMiddleware(logger *utils.Logger) *AuthMiddleware {
	return &AuthMiddleware{logger: logger}
}

// Handler returns the auth middleware handler.
func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For now, this is a placeholder that just logs and passes through.
		// In a real implementation, you'd validate JWT tokens, API keys, etc.

		token := am.extractToken(r)
		if token != "" {
			am.logger.Debug("Authentication token found: %s...", token[:min(len(token), 10)])

			// TODO: Validate token and extract user information.
			// For now, we'll just add a placeholder user to the context.
			ctx := context.WithValue(r.Context(), "user_id", "anonymous")
			ctx = context.WithValue(ctx, "user_role", "user")
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// RequireAuthMiddleware requires authentication for protected routes.
type RequireAuthMiddleware struct {
	logger   *utils.Logger
	response *utils.ResponseHelper
}

// NewRequireAuthMiddleware creates a middleware that requires authentication.
func NewRequireAuthMiddleware(logger *utils.Logger) *RequireAuthMiddleware {
	return &RequireAuthMiddleware{
		logger:   logger,
		response: utils.NewResponseHelper(),
	}
}

// Handler returns the require auth middleware handler.
func (ram *RequireAuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ram.extractToken(r)

		if token == "" {
			ram.logger.Warn("Unauthorized access attempt to %s from %s", r.URL.Path, r.RemoteAddr)
			ram.response.SendError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// TODO: Validate token.
		// For now, we accept any non-empty token.

		ctx := context.WithValue(r.Context(), "user_id", "authenticated_user")
		ctx = context.WithValue(ctx, "user_role", "user")
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// RoleMiddleware checks if user has required role.
type RoleMiddleware struct {
	requiredRole string
	logger       *utils.Logger
	response     *utils.ResponseHelper
}

// NewRoleMiddleware creates a middleware that requires a specific role.
func NewRoleMiddleware(requiredRole string, logger *utils.Logger) *RoleMiddleware {
	return &RoleMiddleware{
		requiredRole: requiredRole,
		logger:       logger,
		response:     utils.NewResponseHelper(),
	}
}

// Handler returns the role middleware handler.
func (rm *RoleMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole, ok := r.Context().Value("user_role").(string)
		if !ok {
			rm.logger.Warn("No user role found in context for %s", r.URL.Path)
			rm.response.SendError(w, http.StatusForbidden, "Access denied")
			return
		}

		if !rm.hasRequiredRole(userRole, rm.requiredRole) {
			rm.logger.Warn("User with role %s attempted to access %s (requires %s)", userRole, r.URL.Path, rm.requiredRole)
			rm.response.SendError(w, http.StatusForbidden, "Insufficient permissions")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper methods.

func (am *AuthMiddleware) extractToken(r *http.Request) string {
	// Check Authorization header.
	auth := r.Header.Get("Authorization")
	if auth != "" {
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Check query parameter.
	if token := r.URL.Query().Get("token"); token != "" {
		return token
	}

	return ""
}

func (ram *RequireAuthMiddleware) extractToken(r *http.Request) string {
	return (&AuthMiddleware{}).extractToken(r)
}

func (rm *RoleMiddleware) hasRequiredRole(userRole, requiredRole string) bool {
	// Simple role hierarchy: admin > user > viewer.
	roleHierarchy := map[string]int{
		"viewer": 1,
		"user":   2,
		"admin":  3,
	}

	userLevel, userExists := roleHierarchy[userRole]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}

// Helper function for minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
