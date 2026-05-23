package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"project-alpha/backend/internal/apperror"
	"project-alpha/backend/internal/dto"
	"project-alpha/backend/internal/models"
	"project-alpha/backend/internal/repositories"
	"project-alpha/backend/internal/validation"
)

type TicketService struct {
	tickets repositories.TicketStore
}

func NewTicketService(tickets repositories.TicketStore) *TicketService {
	return &TicketService{tickets: tickets}
}

func (s *TicketService) Create(req dto.CreateTicketRequest) (models.Ticket, error) {
	title := validation.NormalizeTicketTitle(req.Title)
	description := validation.NormalizeTicketDescription(req.Description)

	if err := validation.ValidateTicket(title, description, req.TagIDs); err != nil {
		return models.Ticket{}, err
	}

	var created models.Ticket
	if err := s.tickets.Transaction(func(repo repositories.TicketStore) error {
		tags, err := s.findAllTags(repo, req.TagIDs)
		if err != nil {
			return err
		}

		ticket := models.Ticket{
			Title:       title,
			Description: description,
			Completed:   false,
			CompletedAt: nil,
		}
		if err := repo.Create(&ticket); err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to create ticket", err)
		}

		if err := repo.ReplaceTags(&ticket, tags); err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to set ticket tags", err)
		}

		created, err = repo.FindByID(ticket.ID)
		if err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to load created ticket", err)
		}

		return nil
	}); err != nil {
		return models.Ticket{}, err
	}

	return created, nil
}

func (s *TicketService) List(query dto.ListTicketsQuery) ([]models.Ticket, dto.PaginationResponse, error) {
	filter, err := validation.ParseListTicketsFilter(query)
	if err != nil {
		return nil, dto.PaginationResponse{}, err
	}

	tickets, total, err := s.tickets.List(filter)
	if err != nil {
		return nil, dto.PaginationResponse{}, apperror.Wrap(apperror.ErrDatabase, "failed to list tickets", err)
	}

	return tickets, dto.NewPaginationResponse(filter.Page, filter.PageSize, total), nil
}

func (s *TicketService) Get(id uint) (models.Ticket, error) {
	ticket, err := s.tickets.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Ticket{}, apperror.New(apperror.ErrNotFound, "ticket not found")
		}
		return models.Ticket{}, apperror.Wrap(apperror.ErrDatabase, "failed to find ticket", err)
	}
	return ticket, nil
}

func (s *TicketService) Update(id uint, req dto.UpdateTicketRequest) (models.Ticket, error) {
	title := validation.NormalizeTicketTitle(req.Title)
	description := validation.NormalizeTicketDescription(req.Description)

	if err := validation.ValidateTicket(title, description, req.TagIDs); err != nil {
		return models.Ticket{}, err
	}

	var updated models.Ticket
	if err := s.tickets.Transaction(func(repo repositories.TicketStore) error {
		ticket, err := repo.FindByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.ErrNotFound, "ticket not found")
			}
			return apperror.Wrap(apperror.ErrDatabase, "failed to find ticket", err)
		}

		tags, err := s.findAllTags(repo, req.TagIDs)
		if err != nil {
			return err
		}

		ticket.Title = title
		ticket.Description = description
		if err := repo.Save(&ticket); err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to update ticket", err)
		}

		if err := repo.ReplaceTags(&ticket, tags); err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to replace ticket tags", err)
		}

		updated, err = repo.FindByID(ticket.ID)
		if err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to load updated ticket", err)
		}

		return nil
	}); err != nil {
		return models.Ticket{}, err
	}

	return updated, nil
}

func (s *TicketService) Delete(id uint) error {
	return s.tickets.Transaction(func(repo repositories.TicketStore) error {
		ticket, err := repo.FindByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.ErrNotFound, "ticket not found")
			}
			return apperror.Wrap(apperror.ErrDatabase, "failed to find ticket", err)
		}

		if err := repo.Delete(&ticket); err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to delete ticket", err)
		}

		return nil
	})
}

func (s *TicketService) Complete(id uint) (models.Ticket, error) {
	return s.setCompletion(id, true)
}

func (s *TicketService) Uncomplete(id uint) (models.Ticket, error) {
	return s.setCompletion(id, false)
}

func (s *TicketService) AddTag(ticketID uint, tagID uint) (models.Ticket, error) {
	var updated models.Ticket
	if err := s.tickets.Transaction(func(repo repositories.TicketStore) error {
		ticket, tag, err := s.findTicketAndTag(repo, ticketID, tagID)
		if err != nil {
			return err
		}

		if err := repo.AddTag(&ticket, tag); err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to add ticket tag", err)
		}

		updated, err = repo.FindByID(ticket.ID)
		if err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to load updated ticket", err)
		}

		return nil
	}); err != nil {
		return models.Ticket{}, err
	}

	return updated, nil
}

func (s *TicketService) RemoveTag(ticketID uint, tagID uint) (models.Ticket, error) {
	var updated models.Ticket
	if err := s.tickets.Transaction(func(repo repositories.TicketStore) error {
		ticket, tag, err := s.findTicketAndTag(repo, ticketID, tagID)
		if err != nil {
			return err
		}

		if err := repo.RemoveTag(&ticket, tag); err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to remove ticket tag", err)
		}

		updated, err = repo.FindByID(ticket.ID)
		if err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to load updated ticket", err)
		}

		return nil
	}); err != nil {
		return models.Ticket{}, err
	}

	return updated, nil
}

func (s *TicketService) setCompletion(id uint, completed bool) (models.Ticket, error) {
	var updated models.Ticket
	if err := s.tickets.Transaction(func(repo repositories.TicketStore) error {
		ticket, err := repo.FindByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.ErrNotFound, "ticket not found")
			}
			return apperror.Wrap(apperror.ErrDatabase, "failed to find ticket", err)
		}

		ticket.Completed = completed
		if completed {
			now := time.Now().UTC()
			ticket.CompletedAt = &now
		} else {
			ticket.CompletedAt = nil
		}

		if err := repo.Save(&ticket); err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to update ticket completion", err)
		}

		updated, err = repo.FindByID(ticket.ID)
		if err != nil {
			return apperror.Wrap(apperror.ErrDatabase, "failed to load updated ticket", err)
		}

		return nil
	}); err != nil {
		return models.Ticket{}, err
	}

	return updated, nil
}

func (s *TicketService) findAllTags(repo repositories.TicketStore, tagIDs []uint) ([]models.Tag, error) {
	tags, err := repo.FindTagsByIDs(tagIDs)
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrDatabase, "failed to find tags", err)
	}

	if len(tags) != len(uniqueUint(tagIDs)) {
		return nil, apperror.New(apperror.ErrNotFound, "one or more tags were not found")
	}

	return tags, nil
}

func (s *TicketService) findTicketAndTag(repo repositories.TicketStore, ticketID uint, tagID uint) (models.Ticket, models.Tag, error) {
	ticket, err := repo.FindByID(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Ticket{}, models.Tag{}, apperror.New(apperror.ErrNotFound, "ticket not found")
		}
		return models.Ticket{}, models.Tag{}, apperror.Wrap(apperror.ErrDatabase, "failed to find ticket", err)
	}

	tag, err := repo.FindTagByID(tagID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Ticket{}, models.Tag{}, apperror.New(apperror.ErrNotFound, "tag not found")
		}
		return models.Ticket{}, models.Tag{}, apperror.Wrap(apperror.ErrDatabase, "failed to find tag", err)
	}

	return ticket, tag, nil
}

func uniqueUint(values []uint) []uint {
	seen := make(map[uint]struct{}, len(values))
	unique := make([]uint, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		unique = append(unique, value)
	}
	return unique
}
