package handlers

import (
	"fmt"
	"net/http"
	"sspr-ldap/domain"
	"sspr-ldap/infra/session"
	"sspr-ldap/infra/template"
	"sspr-ldap/services"
)

type UserHandler struct {
	userService *services.UserService
	session     *session.Store
	template    *template.Engine
}

func NewUserHandler(userService *services.UserService, session *session.Store, template *template.Engine) *UserHandler {
	return &UserHandler{
		userService: userService,
		session:     session,
		template:    template,
	}
}

func (h *UserHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if !h.session.IsAuthenticated(r) {
		fmt.Println("⚠️  Dashboard access denied - not authenticated")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := h.session.GetUsername(r)
	fmt.Printf("✅ Dashboard accessed by: %s\n", username)

	h.template.Render(w, "dashboard.html", map[string]string{
		"Username": username,
	})
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if !h.session.IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := h.session.GetUsername(r)

	if r.Method == "GET" {
		h.template.Render(w, "change-password.html", map[string]string{
			"Username": username,
		})
		return
	}

	if r.Method == "POST" {
		currentPassword := r.FormValue("current_password")
		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")

		if newPassword != confirmPassword {
			w.Header().Set("HX-Retarget", "#error-message")
			w.Header().Set("HX-Reswap", "innerHTML")
			fmt.Fprint(w, `<div class="text-red-600 text-sm">New passwords do not match</div>`)
			return
		}

		userDN := h.session.GetUserDN(r)
		passwordChange := &domain.PasswordChange{
			Username:        username,
			UserDN:          userDN,
			CurrentPassword: currentPassword,
			NewPassword:     newPassword,
		}

		err := h.userService.ChangePassword(passwordChange)
		if err != nil {
			w.Header().Set("HX-Retarget", "#error-message")
			w.Header().Set("HX-Reswap", "innerHTML")
			fmt.Fprintf(w, `<div class="text-red-600 text-sm">Failed to change password: %s</div>`, err.Error())
			return
		}

		w.Header().Set("HX-Retarget", "#success-message")
		w.Header().Set("HX-Reswap", "innerHTML")
		fmt.Fprint(w, `<div class="text-green-600 text-sm">Password changed successfully!</div>`)
		return
	}
}
