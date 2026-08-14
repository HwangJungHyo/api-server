# API Server

A simple Go backend API server providing authentication endpoints (signup, login, profile) and a health check.

## API Endpoints

| Method | Path | Description | Auth Required |
|--------|------|-------------|---------------|
| GET | `/health` | Health check | No |
| POST | `/api/v1/auth/signup` | Register a new user | No |
| POST | `/api/v1/auth/login` | Login | No |
| GET | `/api/v1/auth/me` | Get current user profile | Yes (Bearer token) |

## Quick Start

```bash
# Install dependencies
make init

# Run the server
make run
```

The server starts on port `8080` by default.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `JWT_SECRET` | `dev-secret-change-in-production` | JWT signing secret |

## Development

```bash
# Run tests
make test

# Run tests with coverage
make coverage

# Build binary
make build
```

## Docker

```bash
# Build image
docker build -t api-server .

# Run container
docker run -p 8080:8080 -e JWT_SECRET=your-secret api-server
```

## API Usage Examples

### Sign Up
```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123", "name": "John Doe"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### Get Profile
```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <token>"
```

### Health Check
```bash
curl http://localhost:8080/health
```
