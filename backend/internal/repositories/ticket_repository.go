package repositories

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"project-alpha/backend/internal/dto"
	"project-alpha/backend/internal/models"
)

type TicketRepository struct {
	db *gorm.DB
}

type TicketStore interface {
	Transaction(fn func(repo TicketStore) error) error
	Create(ticket *models.Ticket) error
	List(filter dto.ListTicketsFilter) ([]models.Ticket, int64, error)
	FindByID(id uint) (models.Ticket, error)
	Save(ticket *models.Ticket) error
	Delete(ticket *models.Ticket) error
	ReplaceTags(ticket *models.Ticket, tags []models.Tag) error
	AddTag(ticket *models.Ticket, tag models.Tag) error
	RemoveTag(ticket *models.Ticket, tag models.Tag) error
	FindTagsByIDs(tagIDs []uint) ([]models.Tag, error)
	FindTagByID(id uint) (models.Tag, error)
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Transaction(fn func(repo TicketStore) error) error {
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

func (r *TicketRepository) List(filter dto.ListTicketsFilter) ([]models.Ticket, int64, error) {
	baseQuery := r.applyListFilters(r.db.Model(&models.Ticket{}), filter)

	var total int64
	countQuery := baseQuery.Session(&gorm.Session{})
	if len(filter.TagIDs) > 0 {
		countQuery = countQuery.Select("tickets.id")
	}
	if err := r.db.Table("(?) AS filtered_tickets", countQuery).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count tickets: %w", err)
	}

	var tickets []models.Ticket
	offset := (filter.Page - 1) * filter.PageSize
	if err := baseQuery.Preload("Tags").
		Order("tickets.created_at DESC").
		Limit(filter.PageSize).
		Offset(offset).
		Find(&tickets).Error; err != nil {
		return nil, 0, fmt.Errorf("list tickets: %w", err)
	}

	return tickets, total, nil
}

func (r *TicketRepository) applyListFilters(query *gorm.DB, filter dto.ListTicketsFilter) *gorm.DB {
	if filter.Title != "" {
		query = query.Where("tickets.title ILIKE ?", "%"+filter.Title+"%")
	}

	if filter.Completed != nil {
		query = query.Where("tickets.completed = ?", *filter.Completed)
	}

	if len(filter.TagIDs) > 0 {
		query = query.
			Joins("JOIN ticket_tags ON ticket_tags.ticket_id = tickets.id").
			Where("ticket_tags.tag_id IN ?", filter.TagIDs).
			Group("tickets.id").
			Having("COUNT(DISTINCT ticket_tags.tag_id) = ?", len(filter.TagIDs))
	}

	return query
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
