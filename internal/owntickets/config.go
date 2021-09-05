package owntickets

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func StrToBool(s string) bool {
	switch strings.ToLower(s) {
	case "1":
		return true
	case "yes":
		return true
	case "true":
		return true
	case "y":
		return true
	default:
		return false
	}
}

type Config struct {
	PasswordHash                     string `json:"passwordHash"`
	RequirePasswordForTicketCreation bool   `json:"requirePasswordForTicketCreation"`
	TicketCreationPasswordHash       string `json:"ticketCreationPasswordHash"`
	Database                         string `json:"database"`
}

func (c *Config) Validate() error {
	if len(c.Database) == 0 {
		return fmt.Errorf("set OWNTICKETS_DATABASE environment variable, or set the database key in the config file")
	}
	if len(c.PasswordHash) == 0 {
		return fmt.Errorf("set OWNTICKETS_PASSWORD_HASH environment variable or the passwordHash key in the config file")
	}

	if c.RequirePasswordForTicketCreation && len(c.TicketCreationPasswordHash) == 0 {
		return fmt.Errorf("set OWNTICKETS_TICKET_PASSWORD_HASH environment variable or the ticketCreationPasswordHash key in the config file")
	}

	return nil
}

func (c *Config) LoadFromEnv() {
	passwordHash, ok := os.LookupEnv("OWNTICKETS_PASSWORD_HASH")
	if ok {
		c.PasswordHash = passwordHash
	}

	requirePassword, ok := os.LookupEnv("OWNTICKETS_REQUIRE_PASSWORD")
	if ok {
		c.RequirePasswordForTicketCreation = StrToBool(requirePassword)
	}

	ticketPassword, ok := os.LookupEnv("OWNTICKETS_TICKET_PASSWORD_HASH")
	if ok {
		c.TicketCreationPasswordHash = ticketPassword
	}

	database, ok := os.LookupEnv("OWNTICKETS_DATABASE")
	if ok {
		c.Database = database
	}
}

func (c *Config) LoadFromFile(path string) error {
	var c2 Config

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("couldn't load config from file: %w", err)
	}
	defer f.Close()

	json.NewDecoder(f).Decode(&c2)
	c.PasswordHash = c2.PasswordHash
	c.RequirePasswordForTicketCreation = c2.RequirePasswordForTicketCreation
	c.TicketCreationPasswordHash = c2.TicketCreationPasswordHash
	c.Database = c2.Database

	return nil
}
