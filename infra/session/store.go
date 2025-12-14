package session

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
)

type Store struct {
	store *sessions.CookieStore
}

func NewStore(key string) *Store {
	if len(key) < 32 {
		log.Println("⚠️  WARNING: Session key is too short! Should be at least 32 bytes for security.")
	}

	store := sessions.NewCookieStore([]byte(key))

	// Determine secure cookie setting from environment (default true in production)
	secureDefault := true
	if v := os.Getenv("SESSION_COOKIE_SECURE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			secureDefault = b
		}
	}

	// Configure session options
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
		Secure:   secureDefault,
		SameSite: http.SameSiteLaxMode,
	}

	return &Store{
		store: store,
	}
}

func (s *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	session, err := s.store.Get(r, name)
	if err != nil {
		log.Printf("⚠️  Session Get error: %v", err)
		// Don't return error, return new session instead
		// This handles corrupted cookies gracefully
		session, _ = s.store.New(r, name)
	}
	return session, nil
}

func (s *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	err := session.Save(r, w)
	if err != nil {
		log.Printf("❌ Session Save error: %v", err)
	}
	return err
}

// Helper methods for common session operations
func (s *Store) SetAuthenticated(r *http.Request, w http.ResponseWriter, username, userDN string) error {
	session, err := s.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Values["userDN"] = userDN

	return s.Save(r, w, session)
}

func (s *Store) ClearSession(r *http.Request, w http.ResponseWriter) error {
	session, err := s.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values["authenticated"] = false
	delete(session.Values, "username")
	delete(session.Values, "userDN")

	return s.Save(r, w, session)
}

func (s *Store) IsAuthenticated(r *http.Request) bool {
	session, err := s.Get(r, "session")
	if err != nil {
		return false
	}

	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

func (s *Store) GetUsername(r *http.Request) string {
	session, _ := s.Get(r, "session")
	if username, ok := session.Values["username"].(string); ok {
		return username
	}
	return ""
}

func (s *Store) GetUserDN(r *http.Request) string {
	session, _ := s.Get(r, "session")
	if userDN, ok := session.Values["userDN"].(string); ok {
		return userDN
	}
	return ""
}
