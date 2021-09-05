package main

import (
	log "github.com/sirupsen/logrus"
)

type OwnTickets struct {
	Config Config
	Logger interface{}
}

func (o *OwnTickets) Run() error {
	log.Info("Validating config")
	if err := o.Config.Validate(); err != nil {
		log.WithError(err).Fatal("Config invalid")
		return err
	}
	log.Info("Starting OwnTickets")
	return nil
}
