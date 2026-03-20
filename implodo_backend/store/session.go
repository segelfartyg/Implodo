package store

import (
	"sync"
	"time"
)

// Session holds the temporary state between the app initiating auth
// and it retrieving the JWT after the browser OAuth flow completes.
type Session struct {
	// CodeChallenge is BASE64URL(SHA256(code_verifier)) — provided by the app at start.
	// The app must present the matching code_verifier to retrieve the JWT.
	CodeChallenge string

	// JWT is populated by the Google callback handler after successful auth.
	// Empty string means the auth flow hasn't completed yet.
	JWT string

	CreatedAt time.Time
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	ttl      time.Duration
}

func NewSessionStore() *SessionStore {
	s := &SessionStore{
		sessions: make(map[string]*Session),
		ttl:      10 * time.Minute,
	}
	go s.cleanup()
	return s
}

// Create registers a new pending session keyed by state.
func (s *SessionStore) Create(state, codeChallenge string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[state] = &Session{
		CodeChallenge: codeChallenge,
		CreatedAt:     time.Now(),
	}
}

// SetJWT stores the issued JWT on the session once auth completes.
// Returns false if the state is not found (e.g. expired or invalid).
func (s *SessionStore) SetJWT(state, jwt string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[state]
	if !ok {
		return false
	}
	session.JWT = jwt
	return true
}

// GetAndDelete retrieves the session and removes it — JWT can only be claimed once.
func (s *SessionStore) GetAndDelete(state string) (*Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[state]
	if !ok {
		return nil, false
	}
	delete(s.sessions, state)
	return session, true
}

// cleanup periodically removes expired sessions.
func (s *SessionStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		for state, session := range s.sessions {
			if time.Since(session.CreatedAt) > s.ttl {
				delete(s.sessions, state)
			}
		}
		s.mu.Unlock()
	}
}
