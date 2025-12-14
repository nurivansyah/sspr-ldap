package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sspr-ldap/domain"
	"sspr-ldap/infra/session"
	"sspr-ldap/infra/template"
	"sspr-ldap/services"
)

type AuthHandler struct {
	authService *services.AuthService
	session     *session.Store
	template    *template.Engine
}

func NewAuthHandler(authService *services.AuthService, session *session.Store, template *template.Engine) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		session:     session,
		template:    template,
	}
}

func (h *AuthHandler) Home(w http.ResponseWriter, r *http.Request) {
	if h.session.IsAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Check if already logged in
		if h.session.IsAuthenticated(r) {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
		h.template.Render(w, "login.html", nil)
		return
	}

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		creds := &domain.Credentials{
			Username: username,
			Password: password,
		}

		user, err := h.authService.Authenticate(creds)
		if err != nil {
			log.Printf("Authentication failed: %s", err.Error())
			w.Header().Set("HX-Retarget", "#error-message")
			w.Header().Set("HX-Reswap", "innerHTML")
			fmt.Fprint(w, `<div class="text-red-600 text-sm">Authentication failed</div>`)
			return
		}

		// Set session - CRITICAL: Must be called before any other writes
		err = h.session.SetAuthenticated(r, w, user.Username, user.DN)
		if err != nil {
			log.Printf("Failed to save session: %v\n", err)
			w.Header().Set("HX-Retarget", "#error-message")
			w.Header().Set("HX-Reswap", "innerHTML")
			fmt.Fprint(w, `<div class="text-red-600 text-sm">Failed to create session</div>`)
			return
		}

		log.Printf("session created for user: %s", user.Username)

		// For HTMX requests, use HX-Redirect
		w.Header().Set("HX-Redirect", "/dashboard")
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.session.ClearSession(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
