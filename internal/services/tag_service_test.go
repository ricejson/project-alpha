package services

import (
	"errors"
	"strings"
	"testing"

	"gorm.io/gorm"

	"backend/internal/apperror"
	"backend/internal/dto"
	"backend/internal/models"
)

type fakeTagStore struct {
	nextID uint
	tags   map[uint]models.Tag
}

func newFakeTagStore() *fakeTagStore {
	return &fakeTagStore{nextID: 1, tags: make(map[uint]models.Tag)}
}

func (s *fakeTagStore) Create(tag *models.Tag) error {
	tag.ID = s.nextID
	s.nextID++
	s.tags[tag.ID] = *tag
	return nil
}

func (s *fakeTagStore) List() ([]models.Tag, error) {
	items := make([]models.Tag, 0, len(s.tags))
	for _, tag := range s.tags {
		items = append(items, tag)
	}
	return items, nil
}

func (s *fakeTagStore) FindByID(id uint) (models.Tag, error) {
	tag, ok := s.tags[id]
	if !ok {
		return models.Tag{}, gorm.ErrRecordNotFound
	}
	return tag, nil
}

func (s *fakeTagStore) ExistsByName(name string, excludeID *uint) (bool, error) {
	for _, tag := range s.tags {
		if excludeID != nil && tag.ID == *excludeID {
			continue
		}
		if strings.EqualFold(tag.Name, name) {
			return true, nil
		}
	}
	return false, nil
}

func (s *fakeTagStore) Save(tag *models.Tag) error {
	if _, ok := s.tags[tag.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	s.tags[tag.ID] = *tag
	return nil
}

func (s *fakeTagStore) Delete(tag *models.Tag) error {
	delete(s.tags, tag.ID)
	return nil
}

func TestTagServiceCreateRejectsDuplicateName(t *testing.T) {
	store := newFakeTagStore()
	service := NewTagService(store)

	if _, err := service.Create(dto.CreateTagRequest{Name: "Frontend"}); err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}

	_, err := service.Create(dto.CreateTagRequest{Name: "frontend"})
	if !apperror.Is(err, apperror.ErrConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestTagServiceRejectsInvalidColor(t *testing.T) {
	service := NewTagService(newFakeTagStore())
	color := "red"

	_, err := service.Create(dto.CreateTagRequest{Name: "frontend", Color: &color})
	if !apperror.Is(err, apperror.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestTagServiceDeleteMissingReturnsNotFound(t *testing.T) {
	service := NewTagService(newFakeTagStore())

	err := service.Delete(99)
	if !apperror.Is(err, apperror.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestAppErrorWrapPreservesKind(t *testing.T) {
	err := apperror.Wrap(apperror.ErrDatabase, "database failed", errors.New("boom"))
	if !apperror.Is(err, apperror.ErrDatabase) {
		t.Fatalf("expected database error kind, got %v", err)
	}
}
