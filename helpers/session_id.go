package helpers

import (
	"context"
	"errors"
	"net/http"
)

type SessionIDHelper struct{}

func (s *SessionIDHelper) GetShortenedSessionID(r *http.Request) (string, error) {
	sessionID, ok := r.Context().Value("sessionID).(string)
	if !ok {
		return "", errors.New("session ID not found in context")
	}

	shortenedSessionID := shortenSessionID(sessionID)
	return shortenedSessionID, nil
}

func shortenSessionID(sessionID string) string {
	if len(sessionID) <= 9 {
		return sessionID
	}

	return sessionID[:9]
}

func SetSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, "sessionID", sessionID)
}