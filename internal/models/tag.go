package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name    string   `gorm:"size:40;not null;uniqueIndex"`
	Color   *string  `gorm:"size:7"`
	Tickets []Ticket `gorm:"many2many:ticket_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
