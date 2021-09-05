package owntickets

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/flosch/pongo2/v4"
	"github.com/gorilla/mux"
	"github.com/irth/owntickets/internal/models"
	"gorm.io/gorm"
)

func (o *OwnTickets) ViewTicketPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// TODO: comment submission
		return
	}
	ticketIDs, ok := mux.Vars(r)["id"]
	if !ok {
		// TODO: return 404
		return
	}

	var ticket models.Ticket
	ticketID, err := strconv.Atoi(ticketIDs)
	if err != nil {
		// TODO: err
	}
	err = o.Database.First(&ticket, ticketID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// TODO: return 404
		return
	}

	// TODO: allow admin to access any ticket
	if ticket.Key != r.URL.Query().Get("key") {
		// TODO: return unathorized
		return
	}

	// TODO: display ticket
	o.TicketViewTemplate.ExecuteWriter(pongo2.Context{
		"ticket": ticket,
	}, w)
}
