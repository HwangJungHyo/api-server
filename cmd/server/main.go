package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mentoring-devops/api-server/internal/config"
	"github.com/mentoring-devops/api-server/internal/handler"
	"github.com/mentoring-devops/api-server/internal/repository"
	"github.com/mentoring-devops/api-server/internal/service"
)

const shutdownTimeout = 8 * time.Second


// healthzHandler는 서버가 HTTP 요청을 받을 수 있는지 확인한다.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, "ok")
}

func main() {
	// 1. SIGTERM과 SIGINT를 컨텍스트 취소로 변환
	signalCtx, stopSignal := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stopSignal()

	cfg := config.Load()

	// 의존성 초기화
	userRepo := repository.NewInMemoryUserRepository()
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	h := handler.NewHandler(authService)

	// Setup router (2.기존 애플리케이션 라우터 생성)
	appRouter := handler.NewRouter(h, cfg.JWTSecret)

	// 3. 전체 요청을 받을 상위 라우터 생성
	rootRouter := http.NewServeMux()
	rootRouter.HandleFunc("GET /healthz", healthzHandler)
	
	// 테스트 전용 엔드포인트는 명시적으로 켠 경우에만 등록한다.
	// 운영 이미지에 임의 지연 엔드포인트가 남으면 그 자체가 DoS 표면이 된다.
	if cfg.EnableTestEndpoints {
		log.Println("WARNING: test endpoints enabled (/_test/*) — do not use in production")
		rootRouter.HandleFunc("GET /_test/slow", handler.SlowHandler)
	}

	rootRouter.Handle("/", appRouter)

	addr := fmt.Sprintf(":%s", cfg.Port)

	// 2. Shutdown을 호출할 수 있도록 http.Server 객체를 직접 생성
	srv := &http.Server{
		Addr:              addr,
		Handler:           rootRouter,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// 3. HTTP 서버를 별도 고루틴에서 실행
	go func() {
		log.Printf("Server starting on %s", addr)

		err := srv.ListenAndServe()

		// Shutdown으로 서버가 닫힐 때 발생하는 ErrServerClosed는 정상
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 4. SIGTERM 또는 SIGINT가 올 때까지 메인 고루틴 대기
	<-signalCtx.Done()

	// 두 번째 Ctrl+C는 기본 동작으로 즉시 종료될 수 있도록 복원
	stopSignal()

	log.Println("Shutdown signal received")
	log.Printf("Shutting down server (timeout=%s)...", shutdownTimeout)

	// 5. 최대 8초 동안 진행 중인 요청이 완료되기를 대기
	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		// Shutdown은 시간 초과 시 자동으로 연결을 강제 종료하지 않는다.
		log.Printf("Graceful shutdown failed: %v", err)
		log.Println("Forcing remaining connections to close...")

		if closeErr := srv.Close(); closeErr != nil {
			log.Printf("Forced close failed: %v", closeErr)
		}

	}

	log.Println("Server stopped")
}
