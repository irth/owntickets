package owntickets

import (
	"errors"
	"fmt"
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
		o.Error(w, 404, "Not found", "Ticket ID not provided.")
		return
	}

	var ticket models.Ticket
	ticketID, err := strconv.Atoi(ticketIDs)
	if err != nil {
		o.Error(w, http.StatusUnprocessableEntity, "Invalid ID", "The provided ticket ID is invalid.")
		return
	}
	err = o.Database.First(&ticket, ticketID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		o.Error(w, 404, "Not found", fmt.Sprintf("Ticket #%d could not be found.", ticketID))
		return
	}

	// TODO: allow admin to access any ticket
	if ticket.Key != r.URL.Query().Get("key") {
		o.Error(w, 403, "Unauthorized", "Incorrect access key.")
		return
	}

	// TODO: display ticket
	o.TicketViewTemplate.ExecuteWriter(pongo2.Context{
		"ticket": ticket,
	}, w)
}
