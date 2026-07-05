package models

import "time"

type TicketTag struct {
	TicketID  uint      `gorm:"primaryKey;index"`
	TagID     uint      `gorm:"primaryKey;index"`
	CreatedAt time.Time `gorm:"not null"`
}
