package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mentoring-devops/api-server/internal/config"
	"github.com/mentoring-devops/api-server/internal/handler"
	"github.com/mentoring-devops/api-server/internal/repository"
	"github.com/mentoring-devops/api-server/internal/service"
)

func main() {
	cfg := config.Load()

	// Initialize dependencies
	userRepo := repository.NewInMemoryUserRepository()
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	h := handler.NewHandler(authService)

	// Setup router
	router := handler.NewRouter(h, cfg.JWTSecret)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
