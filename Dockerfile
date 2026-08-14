# Build stage
# go.mod의 go 버전(1.26.5)과 반드시 함께 올릴 것 → troubleshooting/001 참고

FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /server \
    ./cmd/server


# Run stage
FROM alpine:3.24

RUN apk --no-cache add ca-certificates \
    && addgroup -S -g 10001 appgroup \
    && adduser -S -D -H -u 10001 -G appgroup appuser

WORKDIR /app

COPY --from=builder \
    --chown=appuser:appgroup \
    /server ./server

USER appuser:appgroup

EXPOSE 8080

ENTRYPOINT ["./server"]
