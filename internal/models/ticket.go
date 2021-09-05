package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	Name     string
	Email    string
	Title    string
	Content  string
	Priority int
	Key      string
}

func (t *Ticket) CreateKey() error {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	t.Key = uuid.String()
	return nil
}

func (t *Ticket) CheckKey(s string) bool {
	return t.Key == s
}
