package services

import (
	"errors"
	"sort"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"

	"backend/internal/apperror"
	"backend/internal/dto"
	"backend/internal/models"
	"backend/internal/repositories"
)

type fakeTicketStore struct {
	nextTicketID uint
	tickets      map[uint]models.Ticket
	tags         map[uint]models.Tag
	links        map[uint]map[uint]struct{}
}

func newFakeTicketStore() *fakeTicketStore {
	return &fakeTicketStore{
		nextTicketID: 1,
		tickets:      make(map[uint]models.Ticket),
		tags:         make(map[uint]models.Tag),
		links:        make(map[uint]map[uint]struct{}),
	}
}

func (s *fakeTicketStore) withTags(tags ...models.Tag) *fakeTicketStore {
	for _, tag := range tags {
		s.tags[tag.ID] = tag
	}
	return s
}

func (s *fakeTicketStore) Transaction(fn func(repo repositories.TicketStore) error) error {
	return fn(s)
}

func (s *fakeTicketStore) Create(ticket *models.Ticket) error {
	ticket.ID = s.nextTicketID
	s.nextTicketID++
	ticket.CreatedAt = time.Now().UTC()
	ticket.UpdatedAt = ticket.CreatedAt
	s.tickets[ticket.ID] = *ticket
	return nil
}

func (s *fakeTicketStore) List(filter dto.ListTicketsFilter) ([]models.Ticket, int64, error) {
	items := make([]models.Ticket, 0, len(s.tickets))
	for _, ticket := range s.tickets {
		if filter.Title != "" && !strings.Contains(strings.ToLower(ticket.Title), strings.ToLower(filter.Title)) {
			continue
		}
		if filter.Completed != nil && ticket.Completed != *filter.Completed {
			continue
		}
		if len(filter.TagIDs) > 0 && !s.ticketHasAllTags(ticket.ID, filter.TagIDs) {
			continue
		}
		items = append(items, s.attachTags(ticket))
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	total := int64(len(items))
	start := (filter.Page - 1) * filter.PageSize
	if start >= len(items) {
		return []models.Ticket{}, total, nil
	}
	end := start + filter.PageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], total, nil
}

func (s *fakeTicketStore) FindByID(id uint) (models.Ticket, error) {
	ticket, ok := s.tickets[id]
	if !ok {
		return models.Ticket{}, gorm.ErrRecordNotFound
	}
	return s.attachTags(ticket), nil
}

func (s *fakeTicketStore) Save(ticket *models.Ticket) error {
	if _, ok := s.tickets[ticket.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	s.tickets[ticket.ID] = *ticket
	return nil
}

func (s *fakeTicketStore) Delete(ticket *models.Ticket) error {
	delete(s.links, ticket.ID)
	delete(s.tickets, ticket.ID)
	return nil
}

func (s *fakeTicketStore) ReplaceTags(ticket *models.Ticket, tags []models.Tag) error {
	s.links[ticket.ID] = make(map[uint]struct{}, len(tags))
	for _, tag := range tags {
		s.links[ticket.ID][tag.ID] = struct{}{}
	}
	return nil
}

func (s *fakeTicketStore) AddTag(ticket *models.Ticket, tag models.Tag) error {
	if s.links[ticket.ID] == nil {
		s.links[ticket.ID] = make(map[uint]struct{})
	}
	s.links[ticket.ID][tag.ID] = struct{}{}
	return nil
}

func (s *fakeTicketStore) RemoveTag(ticket *models.Ticket, tag models.Tag) error {
	delete(s.links[ticket.ID], tag.ID)
	return nil
}

func (s *fakeTicketStore) FindTagsByIDs(tagIDs []uint) ([]models.Tag, error) {
	tags := make([]models.Tag, 0, len(tagIDs))
	seen := make(map[uint]struct{}, len(tagIDs))
	for _, id := range tagIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		if tag, ok := s.tags[id]; ok {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

func (s *fakeTicketStore) FindTagByID(id uint) (models.Tag, error) {
	tag, ok := s.tags[id]
	if !ok {
		return models.Tag{}, gorm.ErrRecordNotFound
	}
	return tag, nil
}

func (s *fakeTicketStore) attachTags(ticket models.Ticket) models.Ticket {
	ticket.Tags = nil
	for tagID := range s.links[ticket.ID] {
		if tag, ok := s.tags[tagID]; ok {
			ticket.Tags = append(ticket.Tags, tag)
		}
	}
	sort.Slice(ticket.Tags, func(i, j int) bool { return ticket.Tags[i].ID < ticket.Tags[j].ID })
	return ticket
}

func (s *fakeTicketStore) ticketHasAllTags(ticketID uint, tagIDs []uint) bool {
	for _, tagID := range tagIDs {
		if _, ok := s.links[ticketID][tagID]; !ok {
			return false
		}
	}
	return true
}

func TestTicketServiceCreateUpdateAndTagAssociations(t *testing.T) {
	store := newFakeTicketStore().withTags(tagFixture(1, "frontend"), tagFixture(2, "backend"))
	service := NewTicketService(store)
	ticket, err := service.Create(dto.CreateTicketRequest{
		Title:       " Build list ",
		Description: " description ",
		TagIDs:      []uint{1},
	})
	if err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}
	if ticket.Completed || ticket.CompletedAt != nil {
		t.Fatalf("expected new ticket to be incomplete")
	}
	if ticket.Title != "Build list" {
		t.Fatalf("expected trimmed title, got %q", ticket.Title)
	}
	if got := len(ticket.Tags); got != 1 {
		t.Fatalf("expected 1 tag, got %d", got)
	}

	updated, err := service.Update(ticket.ID, dto.UpdateTicketRequest{
		Title:       "Build API",
		Description: "new",
		TagIDs:      []uint{2},
	})
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}
	if got := len(updated.Tags); got != 1 || updated.Tags[0].ID != 2 {
		t.Fatalf("expected tags to be replaced with tag 2, got %+v", updated.Tags)
	}
}

func TestTicketServiceCompleteAndUncomplete(t *testing.T) {
	store := newFakeTicketStore()
	service := NewTicketService(store)
	ticket, err := service.Create(dto.CreateTicketRequest{Title: "task"})
	if err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}

	completed, err := service.Complete(ticket.ID)
	if err != nil {
		t.Fatalf("expected complete to succeed, got %v", err)
	}
	if !completed.Completed || completed.CompletedAt == nil {
		t.Fatalf("expected completed ticket with completedAt")
	}

	uncompleted, err := service.Uncomplete(ticket.ID)
	if err != nil {
		t.Fatalf("expected uncomplete to succeed, got %v", err)
	}
	if uncompleted.Completed || uncompleted.CompletedAt != nil {
		t.Fatalf("expected uncompleted ticket without completedAt")
	}
}

func TestTicketServiceAddRemoveTagAndListFilters(t *testing.T) {
	store := newFakeTicketStore().withTags(tagFixture(1, "frontend"), tagFixture(2, "backend"))
	service := NewTicketService(store)

	first, err := service.Create(dto.CreateTicketRequest{Title: "Frontend task", TagIDs: []uint{1}})
	if err != nil {
		t.Fatalf("create first ticket: %v", err)
	}
	second, err := service.Create(dto.CreateTicketRequest{Title: "Backend task", TagIDs: []uint{1, 2}})
	if err != nil {
		t.Fatalf("create second ticket: %v", err)
	}
	if _, err := service.AddTag(first.ID, 2); err != nil {
		t.Fatalf("add tag: %v", err)
	}
	if _, err := service.AddTag(first.ID, 2); err != nil {
		t.Fatalf("duplicate add tag should be idempotent: %v", err)
	}
	if _, err := service.RemoveTag(first.ID, 2); err != nil {
		t.Fatalf("remove tag: %v", err)
	}

	completed := false
	items, pagination, err := service.List(dto.ListTicketsQuery{
		Title:     "task",
		TagIDsRaw: "1,2",
		Completed: &completed,
		Page:      1,
		PageSize:  10,
	})
	if err != nil {
		t.Fatalf("list tickets: %v", err)
	}
	if pagination.Total != 1 || len(items) != 1 || items[0].ID != second.ID {
		t.Fatalf("expected only second ticket for multi-tag AND filter, total=%d items=%+v", pagination.Total, items)
	}
}

func tagFixture(id uint, name string) models.Tag {
	tag := models.Tag{Name: name}
	tag.ID = id
	return tag
}

func TestTicketServiceMissingResources(t *testing.T) {
	service := NewTicketService(newFakeTicketStore())

	if _, err := service.Get(99); !apperror.Is(err, apperror.ErrNotFound) {
		t.Fatalf("expected missing ticket not found, got %v", err)
	}
	if _, err := service.Create(dto.CreateTicketRequest{Title: "task", TagIDs: []uint{99}}); !apperror.Is(err, apperror.ErrNotFound) {
		t.Fatalf("expected missing tag not found, got %v", err)
	}
	if err := service.Delete(99); !apperror.Is(err, apperror.ErrNotFound) {
		t.Fatalf("expected delete missing ticket not found, got %v", err)
	}
	if _, err := service.Complete(99); !apperror.Is(err, apperror.ErrNotFound) {
		t.Fatalf("expected complete missing ticket not found, got %v", err)
	}
	if _, err := service.AddTag(99, 1); !apperror.Is(err, apperror.ErrNotFound) {
		t.Fatalf("expected add tag missing ticket not found, got %v", err)
	}
	if _, err := service.AddTag(1, 99); !apperror.Is(err, apperror.ErrNotFound) && !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected add tag missing tag not found, got %v", err)
	}
}
