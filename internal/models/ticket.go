package models

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	Title    string
	Content  string
	Priority int
}
