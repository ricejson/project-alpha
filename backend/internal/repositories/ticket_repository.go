package repositories

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"project-alpha/backend/internal/models"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Transaction(fn func(repo *TicketRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewTicketRepository(tx))
	})
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	if err := r.db.Create(ticket).Error; err != nil {
		return fmt.Errorf("create ticket: %w", err)
	}
	return nil
}

func (r *TicketRepository) List(page int, pageSize int) ([]models.Ticket, int64, error) {
	var total int64
	if err := r.db.Model(&models.Ticket{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count tickets: %w", err)
	}

	var tickets []models.Ticket
	offset := (page - 1) * pageSize
	if err := r.db.Preload("Tags").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&tickets).Error; err != nil {
		return nil, 0, fmt.Errorf("list tickets: %w", err)
	}

	return tickets, total, nil
}

func (r *TicketRepository) FindByID(id uint) (models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.Preload("Tags").First(&ticket, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Ticket{}, gorm.ErrRecordNotFound
		}
		return models.Ticket{}, fmt.Errorf("find ticket by id: %w", err)
	}
	return ticket, nil
}

func (r *TicketRepository) Save(ticket *models.Ticket) error {
	if err := r.db.Save(ticket).Error; err != nil {
		return fmt.Errorf("save ticket: %w", err)
	}
	return nil
}

func (r *TicketRepository) Delete(ticket *models.Ticket) error {
	if err := r.db.Model(ticket).Association("Tags").Clear(); err != nil {
		return fmt.Errorf("clear ticket tags: %w", err)
	}

	if err := r.db.Delete(ticket).Error; err != nil {
		return fmt.Errorf("delete ticket: %w", err)
	}
	return nil
}

func (r *TicketRepository) ReplaceTags(ticket *models.Ticket, tags []models.Tag) error {
	if err := r.db.Model(ticket).Association("Tags").Replace(tags); err != nil {
		return fmt.Errorf("replace ticket tags: %w", err)
	}
	return nil
}

func (r *TicketRepository) AddTag(ticket *models.Ticket, tag models.Tag) error {
	var count int64
	if err := r.db.Model(&models.TicketTag{}).
		Where("ticket_id = ? AND tag_id = ?", ticket.ID, tag.ID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("count ticket tag association: %w", err)
	}
	if count > 0 {
		return nil
	}

	if err := r.db.Model(ticket).Association("Tags").Append(&tag); err != nil {
		return fmt.Errorf("add ticket tag: %w", err)
	}
	return nil
}

func (r *TicketRepository) RemoveTag(ticket *models.Ticket, tag models.Tag) error {
	if err := r.db.Model(ticket).Association("Tags").Delete(&tag); err != nil {
		return fmt.Errorf("remove ticket tag: %w", err)
	}
	return nil
}

func (r *TicketRepository) FindTagsByIDs(tagIDs []uint) ([]models.Tag, error) {
	if len(tagIDs) == 0 {
		return []models.Tag{}, nil
	}

	var tags []models.Tag
	if err := r.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("find tags by ids: %w", err)
	}
	return tags, nil
}

func (r *TicketRepository) FindTagByID(id uint) (models.Tag, error) {
	var tag models.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Tag{}, gorm.ErrRecordNotFound
		}
		return models.Tag{}, fmt.Errorf("find tag by id: %w", err)
	}
	return tag, nil
}
