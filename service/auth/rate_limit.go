package auth

import (
	"fmt"
	"time"
)

const (
	maxFailedAttempts  = 5
	attemptWindow      = 15 * time.Minute
	attemptLockoutTime = 15 * time.Minute
)

func (s *authService) checkAttemptAllowed(store map[string]attemptState, key string, now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, ok := store[key]
	if !ok {
		return nil
	}

	if !state.LockedUntil.IsZero() && now.Before(state.LockedUntil) {
		return ErrTooManyAttempts
	}

	if !state.FirstFailed.IsZero() && now.Sub(state.FirstFailed) > attemptWindow {
		delete(store, key)
	}

	return nil
}

func (s *authService) recordAttemptFailure(store map[string]attemptState, key string, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := store[key]
	if state.FirstFailed.IsZero() || now.Sub(state.FirstFailed) > attemptWindow {
		state = attemptState{
			Count:       1,
			FirstFailed: now,
		}
	} else {
		state.Count++
	}

	if state.Count >= maxFailedAttempts {
		state.LockedUntil = now.Add(attemptLockoutTime)
	}

	store[key] = state
}

func (s *authService) clearAttemptState(store map[string]attemptState, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(store, key)
}

func loginAttemptKeyForEmail(email string) string {
	return fmt.Sprintf("login:email:%s", email)
}

func loginAttemptKeyForIP(clientIP string) string {
	return fmt.Sprintf("login:ip:%s", clientIP)
}

func verify2FAAttemptKeyForUser(userID string) string {
	return fmt.Sprintf("verify2fa:user:%s", userID)
}

func verify2FAAttemptKeyForIP(clientIP string) string {
	return fmt.Sprintf("verify2fa:ip:%s", clientIP)
}
