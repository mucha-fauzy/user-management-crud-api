package context

import (
	"context"
	"errors"
	"net/http"
)

const (
	userIDKey   = "user_id"
	usernameKey = "username"
)

// GetUserIDFromContext retrieves the user ID from the request context.
func GetUserIDFromContext(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

// GetUsernameFromContext retrieves the username from the request context.
func GetUsernameFromContext(r *http.Request) (string, error) {
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok {
		return "", errors.New("username not found in context")
	}
	return username, nil
}

// WithUserID adds the user ID to the context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// WithUsername adds the username to the context.
func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey, username)
}
