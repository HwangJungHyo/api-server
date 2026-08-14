package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mentoring-devops/api-server/internal/model"
	"github.com/mentoring-devops/api-server/internal/repository"
	"github.com/mentoring-devops/api-server/internal/service"
)

const testJWTSecret = "test-secret"

func setupTestHandler() *Handler {
	repo := repository.NewInMemoryUserRepository()
	authService := service.NewAuthService(repo, testJWTSecret)
	return NewHandler(authService)
}

func TestHealthCheck(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	h.HealthCheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp model.HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %s", resp.Status)
	}
}

func TestHealthCheck_MethodNotAllowed(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rr := httptest.NewRecorder()

	h.HealthCheck(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestSignup_Success(t *testing.T) {
	h := setupTestHandler()

	body := model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Signup(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}

	var resp model.AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token to be non-empty")
	}
	if resp.User.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", resp.User.Email)
	}
}

func TestSignup_MethodNotAllowed(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/signup", nil)
	rr := httptest.NewRecorder()

	h.Signup(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestSignup_InvalidBody(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Signup(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestSignup_ValidationError(t *testing.T) {
	h := setupTestHandler()

	body := model.SignupRequest{
		Email:    "",
		Password: "password123",
		Name:     "Test User",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Signup(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestSignup_Conflict(t *testing.T) {
	h := setupTestHandler()

	body := model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	bodyBytes, _ := json.Marshal(body)

	// First signup
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Signup(rr, req)

	// Second signup with same email
	bodyBytes, _ = json.Marshal(body)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	h.Signup(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", rr.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	h := setupTestHandler()

	// First signup
	signupBody := model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	bodyBytes, _ := json.Marshal(signupBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Signup(rr, req)

	// Then login
	loginBody := model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	bodyBytes, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	h.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp model.AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token to be non-empty")
	}
}

func TestLogin_MethodNotAllowed(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/login", nil)
	rr := httptest.NewRecorder()

	h.Login(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestLogin_InvalidBody(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Login(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestLogin_ValidationError(t *testing.T) {
	h := setupTestHandler()

	body := model.LoginRequest{
		Email:    "",
		Password: "password123",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Login(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	h := setupTestHandler()

	body := model.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestGetProfile_MethodNotAllowed(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/me", nil)
	rr := httptest.NewRecorder()

	h.GetProfile(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestGetProfile_Unauthorized(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	rr := httptest.NewRecorder()

	h.GetProfile(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	// Simulate auth middleware setting user ID
	ctx := req.Context()
	ctx = setUserIDInContext(ctx, "nonexistent-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	h.GetProfile(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}
