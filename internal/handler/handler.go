package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mentoring-devops/api-server/internal/middleware"
	"github.com/mentoring-devops/api-server/internal/model"
	"github.com/mentoring-devops/api-server/internal/repository"
	"github.com/mentoring-devops/api-server/internal/service"
)

// Handler holds all HTTP handler dependencies.
type Handler struct {
	authService *service.AuthService
}

// NewHandler creates a new Handler instance.
func NewHandler(authService *service.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

// HealthCheck handles GET /health requests.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is allowed")
		return
	}

	resp := model.HealthResponse{
		Status: "ok",
		Time:   time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, resp)
}

// Signup handles POST /api/v1/auth/signup requests.
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only POST is allowed")
		return
	}

	var req model.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	resp, err := h.authService.Signup(&req)
	if err != nil {
		switch err {
		case service.ErrEmailRequired, service.ErrPasswordRequired,
			service.ErrNameRequired, service.ErrPasswordTooShort:
			writeError(w, http.StatusBadRequest, "validation_error", err.Error())
		case repository.ErrUserAlreadyExists:
			writeError(w, http.StatusConflict, "conflict", "user with this email already exists")
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "something went wrong")
		}
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// Login handles POST /api/v1/auth/login requests.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only POST is allowed")
		return
	}

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		switch err {
		case service.ErrEmailRequired, service.ErrPasswordRequired:
			writeError(w, http.StatusBadRequest, "validation_error", err.Error())
		case service.ErrInvalidCredentials:
			writeError(w, http.StatusUnauthorized, "unauthorized", err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "something went wrong")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetProfile handles GET /api/v1/auth/me requests.
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is allowed")
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "user not authenticated")
		return
	}

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		switch err {
		case repository.ErrUserNotFound:
			writeError(w, http.StatusNotFound, "not_found", "user not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "something went wrong")
		}
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, errCode, message string) {
	resp := model.ErrorResponse{
		Error:   errCode,
		Message: message,
	}
	writeJSON(w, status, resp)
}
