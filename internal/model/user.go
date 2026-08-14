package model

import "time"

// User represents a user in the system.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SignupRequest represents the request body for user registration.
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest represents the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the response after successful authentication.
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}
