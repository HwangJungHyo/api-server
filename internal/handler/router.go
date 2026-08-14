package handler

import (
	"net/http"

	"github.com/mentoring-devops/api-server/internal/middleware"
)

// NewRouter sets up and returns the HTTP router.
func NewRouter(h *Handler, jwtSecret string) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/api/v1/auth/signup", h.Signup)
	mux.HandleFunc("/api/v1/auth/login", h.Login)

	// Protected routes
	authMiddleware := middleware.AuthMiddleware(jwtSecret)
	mux.Handle("/api/v1/auth/me", authMiddleware(http.HandlerFunc(h.GetProfile)))

	return mux
}
