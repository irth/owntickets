package owntickets

import (
	"embed"
	"fmt"
	"net/http"
	"strings"

	"github.com/flosch/pongo2/v4"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/irth/owntickets/internal/models"
	loader "github.com/nathan-osman/pongo2-embed-loader"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed templates
var Templates embed.FS

type OwnTickets struct {
	Config   Config
	Database *gorm.DB
	Router   http.Handler

	TicketFormTemplate *pongo2.Template
}

func (o *OwnTickets) Run() error {
	log.Info("Validating config")
	if err := o.Config.Validate(); err != nil {
		log.WithError(err).Fatal("Config invalid")
		return err
	}
	log.Info("Setting up database")
	if err := o.SetupDatabase(); err != nil {
		log.WithError(err).Fatal("Failed to set up database")
		return err
	}
	log.Info("Starting webserver")
	if err := o.SetupRouter(); err != nil {
		log.WithError(err).Fatal("Failed to set up routing")
	}
	if err := o.SetupTemplates(); err != nil {
		log.WithError(err).Fatal("Failed to load templates")
	}
	http.ListenAndServe(":2137", o.Router)
	return nil
}

func (o *OwnTickets) SetupDatabase() error {
	db, err := gorm.Open(sqlite.Open(o.Config.Database), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&models.Ticket{})
	o.Database = db
	return nil
}

func (o *OwnTickets) SetupRouter() error {
	r := mux.NewRouter()
	r.HandleFunc("/", o.TicketPage)
	r.HandleFunc("/admin", o.AdminPage)
	o.Router = r
	return nil
}

func (o *OwnTickets) SetupTemplates() (err error) {
	templateSet := pongo2.NewSet("", &loader.Loader{Content: Templates})
	o.TicketFormTemplate, err = templateSet.FromFile("templates/ticket_form.html")
	return
}

var decoder = schema.NewDecoder()

type TicketForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Title    string `schema:"title"`
	Content  string `schema:"content"`
	Priority int    `schema:"priority"`
	Password string `schema:"password"`
}

func (t *TicketForm) Validate(requirePass bool) map[string]string {
	errors := make(map[string]string)
	notEmpty := func(key string, value string) bool {
		if value == "" {
			errors[key] = fmt.Sprintf("%s required", strings.Title(key))
			return false
		}
		return true
	}
	if requirePass {
		valid := notEmpty("password", t.Password)
		if valid {
			_ = valid
			// TODO: validate password
		}
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

func (o *OwnTickets) TicketPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		o.TicketFormTemplate.ExecuteWriter(pongo2.Context{
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
	errors := form.Validate(o.Config.RequirePasswordForTicketCreation)
	if len(errors) > 0 {
		o.TicketFormTemplate.ExecuteWriter(pongo2.Context{
			"ask_for_password": o.Config.RequirePasswordForTicketCreation,
			"errors":           errors,
			"form":             form,
		}, w)
		return
	}

}

func (o *OwnTickets) AdminPage(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "hi")
}
