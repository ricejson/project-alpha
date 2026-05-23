package dto

import (
	"time"

	"project-alpha/backend/internal/models"
)

type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	TagIDs      []uint `json:"tagIds"`
}

type UpdateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	TagIDs      []uint `json:"tagIds"`
}

type ListTicketsQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type TicketListResponse struct {
	Items      []TicketResponse   `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

type TicketResponse struct {
	ID          uint          `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Completed   bool          `json:"completed"`
	CompletedAt *time.Time    `json:"completedAt"`
	Tags        []TagResponse `json:"tags"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	DeletedAt   *time.Time    `json:"deletedAt,omitempty"`
}

func NewTicketResponse(ticket models.Ticket) TicketResponse {
	var deletedAt *time.Time
	if ticket.DeletedAt.Valid {
		deletedAt = &ticket.DeletedAt.Time
	}

	return TicketResponse{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Completed:   ticket.Completed,
		CompletedAt: ticket.CompletedAt,
		Tags:        NewTagResponses(ticket.Tags),
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

func NewTicketResponses(tickets []models.Ticket) []TicketResponse {
	items := make([]TicketResponse, 0, len(tickets))
	for _, ticket := range tickets {
		items = append(items, NewTicketResponse(ticket))
	}
	return items
}

func NewPaginationResponse(page int, pageSize int, total int64) PaginationResponse {
	totalPages := 0
	if pageSize > 0 && total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
