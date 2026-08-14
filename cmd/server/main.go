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

// healthzHandler는 서버가 HTTP 요청을 받을 수 있는지 확인한다.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, "ok")
}

func main() {
	cfg := config.Load()

	// Initialize dependencies (1.필요한 객체 생성)
	userRepo := repository.NewInMemoryUserRepository()
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	h := handler.NewHandler(authService)

	// Setup router (2.기존 애플리케이션 라우터 생성)
	appRouter := handler.NewRouter(h, cfg.JWTSecret)

	// 3. 전체 요청을 받을 상위 라우터 생성
	rootRouter := http.NewServeMux()

	// 4. 헬스 체크 요청은 인증 없이 처리
	rootRouter.HandleFunc("GET /healthz", healthzHandler)

	// 5. 나머지 요청은 기존 애플리케이션 라우터로 전달
	rootRouter.Handle("/", appRouter)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, rootRouter); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
