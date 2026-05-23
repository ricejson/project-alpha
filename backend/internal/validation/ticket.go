package validation

import (
	"strings"

	"project-alpha/backend/internal/apperror"
	"project-alpha/backend/internal/dto"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

func NormalizeTicketTitle(title string) string {
	return strings.TrimSpace(title)
}

func NormalizeTicketDescription(description string) string {
	return strings.TrimSpace(description)
}

func ValidateTicket(title string, description string, tagIDs []uint) error {
	details := make(map[string]any)

	if title == "" {
		details["title"] = "title is required"
	} else if len([]rune(title)) > 120 {
		details["title"] = "title must be at most 120 characters"
	}

	if len([]rune(description)) > 5000 {
		details["description"] = "description must be at most 5000 characters"
	}

	for index, tagID := range tagIDs {
		if tagID == 0 {
			details["tagIds"] = map[string]any{
				"index": index,
				"error": "tag id must be a positive integer",
			}
			break
		}
	}

	if len(details) > 0 {
		return apperror.WithDetails(apperror.ErrValidation, "invalid ticket input", details)
	}

	return nil
}

func NormalizeListTicketsQuery(query dto.ListTicketsQuery) dto.ListTicketsQuery {
	if query.Page <= 0 {
		query.Page = 1
	}

	if query.PageSize <= 0 {
		query.PageSize = DefaultPageSize
	}
	if query.PageSize > MaxPageSize {
		query.PageSize = MaxPageSize
	}

	return query
}
