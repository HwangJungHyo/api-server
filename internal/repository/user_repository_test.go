package repository

import (
	"fmt"
	"sync"
	"testing"

	"github.com/mentoring-devops/api-server/internal/model"
)

func TestInMemoryUserRepository_Create(t *testing.T) {
	repo := NewInMemoryUserRepository()

	user := &model.User{
		ID:       "test-id-1",
		Email:    "test@example.com",
		Password: "hashed",
		Name:     "Test User",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify timestamps were set
	if user.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestInMemoryUserRepository_Create_DuplicateEmail(t *testing.T) {
	repo := NewInMemoryUserRepository()

	user1 := &model.User{
		ID:       "test-id-1",
		Email:    "test@example.com",
		Password: "hashed",
		Name:     "User 1",
	}
	user2 := &model.User{
		ID:       "test-id-2",
		Email:    "test@example.com",
		Password: "hashed",
		Name:     "User 2",
	}

	_ = repo.Create(user1)
	err := repo.Create(user2)
	if err != ErrUserAlreadyExists {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestInMemoryUserRepository_FindByEmail(t *testing.T) {
	repo := NewInMemoryUserRepository()

	user := &model.User{
		ID:       "test-id-1",
		Email:    "test@example.com",
		Password: "hashed",
		Name:     "Test User",
	}
	_ = repo.Create(user)

	found, err := repo.FindByEmail("test@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.ID != "test-id-1" {
		t.Errorf("expected ID test-id-1, got %s", found.ID)
	}
}

func TestInMemoryUserRepository_FindByEmail_NotFound(t *testing.T) {
	repo := NewInMemoryUserRepository()

	_, err := repo.FindByEmail("notfound@example.com")
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestInMemoryUserRepository_FindByID(t *testing.T) {
	repo := NewInMemoryUserRepository()

	user := &model.User{
		ID:       "test-id-1",
		Email:    "test@example.com",
		Password: "hashed",
		Name:     "Test User",
	}
	_ = repo.Create(user)

	found, err := repo.FindByID("test-id-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", found.Email)
	}
}

func TestInMemoryUserRepository_FindByID_NotFound(t *testing.T) {
	repo := NewInMemoryUserRepository()

	_, err := repo.FindByID("nonexistent")
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestInMemoryUserRepository_ConcurrentAccess(t *testing.T) {
	repo := NewInMemoryUserRepository()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			user := &model.User{
				ID:       fmt.Sprintf("id-%d", idx),
				Email:    fmt.Sprintf("user%d@example.com", idx),
				Password: "hashed",
				Name:     fmt.Sprintf("User %d", idx),
			}
			_ = repo.Create(user)
		}(i)
	}
	wg.Wait()
}
