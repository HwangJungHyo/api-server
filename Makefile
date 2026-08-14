.PHONY: init run test coverage build clean

# Initialize the project (download dependencies)
init:
	go mod download
	go mod tidy

# Run the server
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./... -v -count=1

# Run tests with coverage
coverage:
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -func=coverage.out

# Build binary
build:
	go build -o bin/server cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out
