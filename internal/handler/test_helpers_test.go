package handler

import (
	"context"

	"github.com/mentoring-devops/api-server/internal/middleware"
)

// setUserIDInContext is a test helper that sets user ID in context.
func setUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, middleware.UserIDKey, userID)
}
