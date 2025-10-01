package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a registered account with per-user content.
type User struct {
	Username     string    `json:"username"`
	PasswordHash []byte    `json:"-"`
	Content      []string  `json:"content"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Store holds users and active sessions in memory.
type Store struct {
	mu       sync.RWMutex
	users    map[string]*User
	sessions map[string]string // token -> username
}

var (
	// ErrUserExists signals the username is already taken.
	ErrUserExists = errors.New("user already exists")
	// ErrInvalidCredentials indicates a login failure.
	ErrInvalidCredentials = errors.New("invalid username or password")
	// ErrUnauthorized indicates an invalid or missing session.
	ErrUnauthorized = errors.New("unauthorized")

	dummyPasswordHash = []byte("$2a$10$CwTycUXWue0Thq9StjUM0uJ8n4VWeNseyX2fA9DE.D7su7J6iYGTC")
)

// NewStore sets up an empty in-memory store.
func NewStore() *Store {
	return &Store{
		users:    make(map[string]*User),
		sessions: make(map[string]string),
	}
}

// CreateUser registers a new user with a hashed password.
func (s *Store) CreateUser(username, password string, content []string) error {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return fmt.Errorf("username and password are required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[username]; exists {
		return ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	s.users[username] = &User{
		Username:     username,
		PasswordHash: hash,
		Content:      append([]string(nil), content...),
		CreatedAt:    time.Now().UTC(),
	}

	return nil
}

// Authenticate validates credentials and returns a session token.
func (s *Store) Authenticate(username, password string) (string, error) {
	s.mu.RLock()
	user, exists := s.users[username]
	s.mu.RUnlock()
	if !exists {
		_ = bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(password))
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := newToken()
	if err != nil {
		return "", fmt.Errorf("create token: %w", err)
	}

	s.mu.Lock()
	s.sessions[token] = user.Username
	s.mu.Unlock()

	return token, nil
}

// ContentByToken returns user-specific content for a valid token.
func (s *Store) ContentByToken(token string) ([]string, error) {
	s.mu.RLock()
	username, ok := s.sessions[token]
	if !ok {
		s.mu.RUnlock()
		return nil, ErrUnauthorized
	}
	user := s.users[username]
	content := append([]string(nil), user.Content...)
	s.mu.RUnlock()
	return content, nil
}

// UpdateContentByToken replaces the content owned by the authenticated user.
func (s *Store) UpdateContentByToken(token string, content []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	username, ok := s.sessions[token]
	if !ok {
		return ErrUnauthorized
	}

	user := s.users[username]
	user.Content = append([]string(nil), content...)
	return nil
}

func newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Server wires HTTP handlers to a Store.
type Server struct {
	store *Store
}

// NewServer configures a Server with the given Store.
func NewServer(store *Store) *Server {
	return &Server{store: store}
}

// Routes exposes the HTTP handlers for account and content management.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/signup", s.handleSignup)
	mux.HandleFunc("/login", s.handleLogin)
	mux.HandleFunc("/me/content", s.handleContent)
	return mux
}

type signupRequest struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Content  []string `json:"content"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON payload"})
		return
	}

	if err := s.store.CreateUser(req.Username, req.Password, req.Content); err != nil {
		switch {
		case errors.Is(err, ErrUserExists):
			writeJSON(w, http.StatusConflict, errorResponse{Error: "username already taken"})
		default:
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON payload"})
		return
	}

	token, err := s.store.Authenticate(req.Username, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		if !errors.Is(err, ErrInvalidCredentials) {
			status = http.StatusInternalServerError
		}
		writeJSON(w, status, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{Token: token})
}

func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	token := parseBearerToken(r.Header.Get("Authorization"))
	if token == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "missing bearer token"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		content, err := s.store.ContentByToken(token)
		if err != nil {
			status := http.StatusUnauthorized
			if !errors.Is(err, ErrUnauthorized) {
				status = http.StatusInternalServerError
			}
			writeJSON(w, status, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, struct {
			Content []string `json:"content"`
		}{Content: content})
	case http.MethodPut:
		var body struct {
			Content []string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON payload"})
			return
		}
		if err := s.store.UpdateContentByToken(token, body.Content); err != nil {
			status := http.StatusUnauthorized
			if !errors.Is(err, ErrUnauthorized) {
				status = http.StatusInternalServerError
			}
			writeJSON(w, status, errorResponse{Error: err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func parseBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}
