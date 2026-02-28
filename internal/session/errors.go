package session

import "errors"

var (
	ErrSessionNotFound  = errors.New("session not found")
	ErrSessionExpired   = errors.New("session expired")
	ErrDuplicateSession = errors.New("duplicate session")
)
