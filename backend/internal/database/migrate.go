package database

import (
	"fmt"

	"gorm.io/gorm"

	"project-alpha/backend/internal/models"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.SetupJoinTable(&models.Ticket{}, "Tags", &models.TicketTag{}); err != nil {
		return fmt.Errorf("setup ticket tags join table: %w", err)
	}

	if err := db.SetupJoinTable(&models.Tag{}, "Tickets", &models.TicketTag{}); err != nil {
		return fmt.Errorf("setup tag tickets join table: %w", err)
	}

	if err := db.AutoMigrate(&models.Tag{}, &models.Ticket{}, &models.TicketTag{}); err != nil {
		return fmt.Errorf("auto migrate database schema: %w", err)
	}

	if err := createAdditionalIndexes(db); err != nil {
		return err
	}

	return nil
}

func createAdditionalIndexes(db *gorm.DB) error {
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_tickets_created_at ON tickets (created_at)").Error; err != nil {
		return fmt.Errorf("create tickets created_at index: %w", err)
	}

	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_lower_name ON tags (lower(name)) WHERE deleted_at IS NULL").Error; err != nil {
		return fmt.Errorf("create tags lower name unique index: %w", err)
	}

	return nil
}
