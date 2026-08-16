package handler

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// maxSlowDuration은 테스트 엔드포인트가 허용하는 지연 상한이다.
const maxSlowDuration = 30 * time.Second

// inFlightSlow는 현재 처리 중인 /_test/slow 요청 수다.
// hey는 서버의 in-flight 요청 수를 알려주지 않으므로 서버가 직접 기록한다.
var inFlightSlow atomic.Int64

// SlowHandler는 지정한 시간만큼 지연 후 200을 반환한다.
// graceful shutdown의 in-flight 요청 보존 실험 전용이며,
// ENABLE_TEST_ENDPOINTS=true일 때만 라우터에 등록된다.
func SlowHandler(w http.ResponseWriter, r *http.Request) {
	duration := 5 * time.Second
	if raw := r.URL.Query().Get("duration"); raw != "" {
		parsed, err := time.ParseDuration(raw)
		if err != nil || parsed <= 0 || parsed > maxSlowDuration {
			http.Error(w, "invalid duration", http.StatusBadRequest)
			return
		}
		duration = parsed
	}

	// 진입 로그. 실험 스크립트는 이 줄을 세어 신호 시점을 결정한다.
	n := inFlightSlow.Add(1)
	log.Printf("slow: enter duration=%s in-flight=%d", duration, n)

	defer func() {
		remaining := inFlightSlow.Add(-1)
		log.Printf("slow: exit in-flight=%d", remaining)
	}()

	select {
	case <-time.After(duration):
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("slow ok\n"))
	case <-r.Context().Done():
		// 클라이언트 연결이 끊기거나 서버가 강제 종료된 경우
		log.Printf("slow: aborted err=%v", r.Context().Err())
	}
}
