#!/usr/bin/env bash
# experiments/002 — graceful shutdown in-flight 요청 보존 실험
# usage: ./exp002-shutdown.sh none|stop|kill
set -euo pipefail

MODE="${1:?usage: $0 none|stop|kill}"
CONTAINER=api-test
IMAGE=api-server:exp002
URL="http://127.0.0.1:8080/_test/slow?duration=5s"
EXPECTED=20

RUN_DIR="results/$(date +%Y%m%d-%H%M%S)-${MODE}"
mkdir -p "$RUN_DIR"

docker rm -f "$CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$CONTAINER" \
  -p 127.0.0.1:8080:8080 \
  -e ENABLE_TEST_ENDPOINTS=true \
  "$IMAGE" >/dev/null

# 동일성 증명: 두 조건이 같은 image ID를 쓰는지 증적에 남긴다
docker image inspect --format '{{.Id}}' "$IMAGE" > "$RUN_DIR/image-id.txt"

# readiness — sleep 대신 실제 응답을 기다린다
until curl -sf http://127.0.0.1:8080/healthz >/dev/null 2>&1; do sleep 0.2; done

hey -n 20 -c 20 -t 15 "$URL" > "$RUN_DIR/hey.txt" 2>&1 &
HEY_PID=$!

# 20개 진입 로그 확인 → 신호. 사람 반응속도를 변수에서 제거한다.
deadline=$((SECONDS + 10))
until [ "$(docker logs "$CONTAINER" 2>&1 | grep -c 'slow: enter')" -ge "$EXPECTED" ]; do
  if [ "$SECONDS" -ge "$deadline" ]; then
    echo "INVALID: ${EXPECTED}개 진입 미달 — 이 실행은 무효" | tee "$RUN_DIR/INVALID.txt"
    kill "$HEY_PID" 2>/dev/null || true
    exit 1
  fi
  sleep 0.05
done
date --rfc-3339=ns > "$RUN_DIR/signal-at.txt"

case "$MODE" in
  none) : ;;                                  # 기준선: 종료 없음
  stop) docker stop --timeout 15 "$CONTAINER" ;;
  kill) docker kill "$CONTAINER" ;;
  *) echo "unknown mode: $MODE" >&2; exit 2 ;;
esac

wait "$HEY_PID" || true
docker logs "$CONTAINER" > "$RUN_DIR/container.log" 2>&1
echo "=== $RUN_DIR ==="
cat "$RUN_DIR/hey.txt"
