package services

import (
	"errors"

	"gorm.io/gorm"

	"project-alpha/backend/internal/apperror"
	"project-alpha/backend/internal/dto"
	"project-alpha/backend/internal/models"
	"project-alpha/backend/internal/repositories"
	"project-alpha/backend/internal/validation"
)

type TagService struct {
	tags repositories.TagStore
}

func NewTagService(tags repositories.TagStore) *TagService {
	return &TagService{tags: tags}
}

func (s *TagService) Create(req dto.CreateTagRequest) (models.Tag, error) {
	name := validation.NormalizeTagName(req.Name)
	color := validation.NormalizeTagColor(req.Color)

	if err := validation.ValidateTag(name, color); err != nil {
		return models.Tag{}, err
	}

	exists, err := s.tags.ExistsByName(name, nil)
	if err != nil {
		return models.Tag{}, apperror.Wrap(apperror.ErrDatabase, "failed to check tag name", err)
	}
	if exists {
		return models.Tag{}, apperror.New(apperror.ErrConflict, "tag name already exists")
	}

	tag := models.Tag{
		Name:  name,
		Color: color,
	}
	if err := s.tags.Create(&tag); err != nil {
		return models.Tag{}, apperror.Wrap(apperror.ErrDatabase, "failed to create tag", err)
	}

	return tag, nil
}

func (s *TagService) List() ([]models.Tag, error) {
	tags, err := s.tags.List()
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrDatabase, "failed to list tags", err)
	}
	return tags, nil
}

func (s *TagService) Update(id uint, req dto.UpdateTagRequest) (models.Tag, error) {
	name := validation.NormalizeTagName(req.Name)
	color := validation.NormalizeTagColor(req.Color)

	if err := validation.ValidateTag(name, color); err != nil {
		return models.Tag{}, err
	}

	tag, err := s.tags.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Tag{}, apperror.New(apperror.ErrNotFound, "tag not found")
		}
		return models.Tag{}, apperror.Wrap(apperror.ErrDatabase, "failed to find tag", err)
	}

	exists, err := s.tags.ExistsByName(name, &id)
	if err != nil {
		return models.Tag{}, apperror.Wrap(apperror.ErrDatabase, "failed to check tag name", err)
	}
	if exists {
		return models.Tag{}, apperror.New(apperror.ErrConflict, "tag name already exists")
	}

	tag.Name = name
	tag.Color = color

	if err := s.tags.Save(&tag); err != nil {
		return models.Tag{}, apperror.Wrap(apperror.ErrDatabase, "failed to update tag", err)
	}

	return tag, nil
}

func (s *TagService) Delete(id uint) error {
	tag, err := s.tags.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(apperror.ErrNotFound, "tag not found")
		}
		return apperror.Wrap(apperror.ErrDatabase, "failed to find tag", err)
	}

	if err := s.tags.Delete(&tag); err != nil {
		return apperror.Wrap(apperror.ErrDatabase, "failed to delete tag", err)
	}

	return nil
}
