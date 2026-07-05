package dto

import (
	"time"

	"backend/internal/models"
)

type CreateTagRequest struct {
	Name  string  `json:"name" binding:"required"`
	Color *string `json:"color"`
}

type UpdateTagRequest struct {
	Name  string  `json:"name" binding:"required"`
	Color *string `json:"color"`
}

type TagResponse struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Color     *string    `json:"color"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func NewTagResponse(tag models.Tag) TagResponse {
	var deletedAt *time.Time
	if tag.DeletedAt.Valid {
		deletedAt = &tag.DeletedAt.Time
	}

	return TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		Color:     tag.Color,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func NewTagResponses(tags []models.Tag) []TagResponse {
	items := make([]TagResponse, 0, len(tags))
	for _, tag := range tags {
		items = append(items, NewTagResponse(tag))
	}
	return items
}
