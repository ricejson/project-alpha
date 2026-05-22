package models

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	Title       string     `gorm:"size:120;not null;index"`
	Description string     `gorm:"type:text;not null;default:''"`
	Completed   bool       `gorm:"not null;default:false;index"`
	CompletedAt *time.Time `gorm:"type:timestamptz"`
	Tags        []Tag      `gorm:"many2many:ticket_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
