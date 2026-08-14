package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/mentoring-devops/api-server/internal/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
}

// InMemoryUserRepository is an in-memory implementation of UserRepository.
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*model.User // keyed by ID
}

// NewInMemoryUserRepository creates a new in-memory user repository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*model.User),
	}
}

// Create stores a new user.
func (r *InMemoryUserRepository) Create(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email already exists
	for _, u := range r.users {
		if u.Email == user.Email {
			return ErrUserAlreadyExists
		}
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	r.users[user.ID] = user
	return nil
}

// FindByEmail finds a user by email address.
func (r *InMemoryUserRepository) FindByEmail(email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

// FindByID finds a user by ID.
func (r *InMemoryUserRepository) FindByID(id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}
