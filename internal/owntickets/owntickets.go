package owntickets

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/gorilla/mux"
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

func (o *OwnTickets) TicketPage(w http.ResponseWriter, r *http.Request) {
	o.TicketFormTemplate.ExecuteWriter(nil, w)
	// TODO: handle POST
}

func (o *OwnTickets) AdminPage(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "hi")
}
