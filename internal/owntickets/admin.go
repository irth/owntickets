package owntickets

import (
	"net/http"

	"github.com/flosch/pongo2/v4"
	"golang.org/x/crypto/bcrypt"
)

func (o *OwnTickets) CheckAdminCookie(w http.ResponseWriter, r *http.Request, require bool) (isAdmin bool) {
	defer func() {
		if require && !isAdmin {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}()

	cookie, err := r.Cookie("admin")
	if err != nil {
		return false
	}
	value := make(map[string]interface{})
	if err = o.Cookie.Decode(cookie.Name, cookie.Value, &value); err == nil {
		isAdmin, ok := value["isAdmin"]
		if !ok {
			return false
		}
		isAdminBool, ok := isAdmin.(bool)
		return ok && isAdminBool
	}
	return false
}

func (o *OwnTickets) SetAdminCookie(w http.ResponseWriter, r *http.Request, isAdmin bool) {
	value := map[string]interface{}{
		"isAdmin": isAdmin,
	}
	if encoded, err := o.Cookie.Encode("admin", value); err == nil {
		cookie := &http.Cookie{
			Name:     "admin",
			Value:    encoded,
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
}

func (o *OwnTickets) AdminPage(w http.ResponseWriter, r *http.Request) {
	if !o.CheckAdminCookie(w, r, true) {
		return
	}
	o.AdminTemplate.ExecuteWriter(pongo2.Context{}, w)
}

func (o *OwnTickets) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		if o.CheckAdminCookie(w, r, false) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
		o.LoginTemplate.ExecuteWriter(pongo2.Context{}, w)
		return
	}

	r.ParseForm()
	pass := r.Form.Get("password")
	if err := bcrypt.CompareHashAndPassword([]byte(o.Config.PasswordHash), []byte(pass)); err != nil {
		o.LoginTemplate.ExecuteWriter(pongo2.Context{"error": true}, w)
		return
	}

	o.SetAdminCookie(w, r, true)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (o *OwnTickets) LogoutPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		o.Error(w, http.StatusMethodNotAllowed, "Method not allowed", "Please try again.")
		return
	}
	o.SetAdminCookie(w, r, false)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
