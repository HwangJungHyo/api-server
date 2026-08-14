package service

import (
	"testing"

	"github.com/mentoring-devops/api-server/internal/model"
	"github.com/mentoring-devops/api-server/internal/repository"
)

const testJWTSecret = "test-secret"

func newTestAuthService() *AuthService {
	repo := repository.NewInMemoryUserRepository()
	return NewAuthService(repo, testJWTSecret)
}

func TestAuthService_Signup_Success(t *testing.T) {
	svc := newTestAuthService()

	req := &model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	resp, err := svc.Signup(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token to be non-empty")
	}
	if resp.User.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", resp.User.Email)
	}
	if resp.User.Name != "Test User" {
		t.Errorf("expected name Test User, got %s", resp.User.Name)
	}
}

func TestAuthService_Signup_EmailRequired(t *testing.T) {
	svc := newTestAuthService()

	req := &model.SignupRequest{
		Email:    "",
		Password: "password123",
		Name:     "Test User",
	}

	_, err := svc.Signup(req)
	if err != ErrEmailRequired {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}

func TestAuthService_Signup_PasswordRequired(t *testing.T) {
	svc := newTestAuthService()

	req := &model.SignupRequest{
		Email:    "test@example.com",
		Password: "",
		Name:     "Test User",
	}

	_, err := svc.Signup(req)
	if err != ErrPasswordRequired {
		t.Fatalf("expected ErrPasswordRequired, got %v", err)
	}
}

func TestAuthService_Signup_PasswordTooShort(t *testing.T) {
	svc := newTestAuthService()

	req := &model.SignupRequest{
		Email:    "test@example.com",
		Password: "12345",
		Name:     "Test User",
	}

	_, err := svc.Signup(req)
	if err != ErrPasswordTooShort {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestAuthService_Signup_NameRequired(t *testing.T) {
	svc := newTestAuthService()

	req := &model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "",
	}

	_, err := svc.Signup(req)
	if err != ErrNameRequired {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestAuthService_Signup_DuplicateEmail(t *testing.T) {
	svc := newTestAuthService()

	req := &model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	_, _ = svc.Signup(req)
	_, err := svc.Signup(req)
	if err != repository.ErrUserAlreadyExists {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	svc := newTestAuthService()

	// First signup
	signupReq := &model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	_, _ = svc.Signup(signupReq)

	// Then login
	loginReq := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := svc.Login(loginReq)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token to be non-empty")
	}
	if resp.User.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", resp.User.Email)
	}
}

func TestAuthService_Login_EmailRequired(t *testing.T) {
	svc := newTestAuthService()

	req := &model.LoginRequest{
		Email:    "",
		Password: "password123",
	}

	_, err := svc.Login(req)
	if err != ErrEmailRequired {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}

func TestAuthService_Login_PasswordRequired(t *testing.T) {
	svc := newTestAuthService()

	req := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "",
	}

	_, err := svc.Login(req)
	if err != ErrPasswordRequired {
		t.Fatalf("expected ErrPasswordRequired, got %v", err)
	}
}

func TestAuthService_Login_InvalidEmail(t *testing.T) {
	svc := newTestAuthService()

	req := &model.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, err := svc.Login(req)
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc := newTestAuthService()

	// First signup
	signupReq := &model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	_, _ = svc.Signup(signupReq)

	// Login with wrong password
	loginReq := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err := svc.Login(loginReq)
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_GetProfile_Success(t *testing.T) {
	svc := newTestAuthService()

	// Signup to create a user
	signupReq := &model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	resp, _ := svc.Signup(signupReq)

	// Get profile
	user, err := svc.GetProfile(resp.User.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}
}

func TestAuthService_GetProfile_NotFound(t *testing.T) {
	svc := newTestAuthService()

	_, err := svc.GetProfile("nonexistent-id")
	if err != repository.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
