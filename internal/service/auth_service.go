package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/mentoring-devops/api-server/internal/model"
	"github.com/mentoring-devops/api-server/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrNameRequired       = errors.New("name is required")
	ErrPasswordTooShort   = errors.New("password must be at least 6 characters")
)

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Signup creates a new user account.
func (s *AuthService) Signup(req *model.SignupRequest) (*model.AuthResponse, error) {
	if err := s.validateSignupRequest(req); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Login authenticates a user and returns a token.
func (s *AuthService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	if req.Email == "" {
		return nil, ErrEmailRequired
	}
	if req.Password == "" {
		return nil, ErrPasswordRequired
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// GetProfile retrieves a user's profile by ID.
func (s *AuthService) GetProfile(userID string) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) validateSignupRequest(req *model.SignupRequest) error {
	if req.Email == "" {
		return ErrEmailRequired
	}
	if req.Password == "" {
		return ErrPasswordRequired
	}
	if len(req.Password) < 6 {
		return ErrPasswordTooShort
	}
	if req.Name == "" {
		return ErrNameRequired
	}
	return nil
}

func (s *AuthService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
