package repositories

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"project-alpha/backend/internal/models"
)

type TagRepository struct {
	db *gorm.DB
}

type TagStore interface {
	Create(tag *models.Tag) error
	List() ([]models.Tag, error)
	FindByID(id uint) (models.Tag, error)
	ExistsByName(name string, excludeID *uint) (bool, error)
	Save(tag *models.Tag) error
	Delete(tag *models.Tag) error
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(tag *models.Tag) error {
	if err := r.db.Create(tag).Error; err != nil {
		return fmt.Errorf("create tag: %w", err)
	}
	return nil
}

func (r *TagRepository) List() ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.Order("name ASC").Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	return tags, nil
}

func (r *TagRepository) FindByID(id uint) (models.Tag, error) {
	var tag models.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Tag{}, gorm.ErrRecordNotFound
		}
		return models.Tag{}, fmt.Errorf("find tag by id: %w", err)
	}
	return tag, nil
}

func (r *TagRepository) ExistsByName(name string, excludeID *uint) (bool, error) {
	query := r.db.Model(&models.Tag{}).Where("lower(name) = lower(?)", name)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check tag name exists: %w", err)
	}

	return count > 0, nil
}

func (r *TagRepository) Save(tag *models.Tag) error {
	if err := r.db.Save(tag).Error; err != nil {
		return fmt.Errorf("save tag: %w", err)
	}
	return nil
}

func (r *TagRepository) Delete(tag *models.Tag) error {
	if err := r.db.Delete(tag).Error; err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	return nil
}
