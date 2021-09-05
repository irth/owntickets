package models

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	Name     string
	Email    string
	Title    string
	Content  string
	Priority int
}
