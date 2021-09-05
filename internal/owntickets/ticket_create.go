package owntickets

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/flosch/pongo2/v4"
	"github.com/irth/owntickets/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type TicketForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Title    string `schema:"title"`
	Content  string `schema:"content"`
	Priority int    `schema:"priority"`
	Password string `schema:"password"`
}

func (t *TicketForm) Validate() map[string]string {
	errors := make(map[string]string)
	notEmpty := func(key string, value string) bool {
		if value == "" {
			errors[key] = fmt.Sprintf("%s required", strings.Title(key))
			return false
		}
		return true
	}
	notEmpty("name", t.Name)
	notEmpty("email", t.Email)
	notEmpty("title", t.Title)
	notEmpty("content", t.Content)

	if t.Priority < 10 {
		errors["priority"] = "Priority too low"
	}
	if t.Priority > 50 {
		errors["priority"] = "Priority too high"
	}

	return errors
}

func (o *OwnTickets) CreateTicketPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		o.TicketCreateTemplate.ExecuteWriter(pongo2.Context{
			"ask_for_password": o.Config.RequirePasswordForTicketCreation,
			"errors":           map[string]string{},
			"form":             TicketForm{Priority: 10},
		}, w)
		return
	}

	var form TicketForm
	r.ParseForm()
	// TODO: handle error from ParseForm
	decoder.Decode(&form, r.Form)
	errors := form.Validate()
	if o.Config.RequirePasswordForTicketCreation {
		err := bcrypt.CompareHashAndPassword([]byte(o.Config.TicketCreationPasswordHash), []byte(form.Password))
		if err != nil {
			errors["password"] = "Incorrect password"
		}
	}

	if len(errors) > 0 {
		o.TicketCreateTemplate.ExecuteWriter(pongo2.Context{
			"ask_for_password": o.Config.RequirePasswordForTicketCreation,
			"errors":           errors,
			"form":             form,
		}, w)
		return
	}

	ticket := models.Ticket{
		Name:     form.Name,
		Email:    form.Email,
		Title:    form.Title,
		Content:  form.Content,
		Priority: form.Priority,
	}
	ticket.CreateKey()
	o.Database.Create(&ticket)
	o.Database.Commit()

	http.Redirect(w, r, fmt.Sprintf("/tickets/%d?key=%s", ticket.ID, ticket.Key), http.StatusSeeOther)
}
