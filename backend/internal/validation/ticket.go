package validation

import (
	"strconv"
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
	query.Title = strings.TrimSpace(query.Title)
	query.TagIDsRaw = strings.TrimSpace(query.TagIDsRaw)

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

func ParseListTicketsFilter(query dto.ListTicketsQuery) (dto.ListTicketsFilter, error) {
	query = NormalizeListTicketsQuery(query)

	tagIDs, err := parseTagIDs(query.TagIDsRaw)
	if err != nil {
		return dto.ListTicketsFilter{}, err
	}

	return dto.ListTicketsFilter{
		Title:     query.Title,
		TagIDs:    tagIDs,
		Completed: query.Completed,
		Page:      query.Page,
		PageSize:  query.PageSize,
	}, nil
}

func parseTagIDs(raw string) ([]uint, error) {
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	tagIDs := make([]uint, 0, len(parts))
	details := make(map[string]any)
	seen := make(map[uint]struct{}, len(parts))

	for index, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			details["tagIds"] = map[string]any{
				"index": index,
				"error": "tag id must not be empty",
			}
			break
		}

		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil || parsed == 0 {
			details["tagIds"] = map[string]any{
				"index": index,
				"error": "tag id must be a positive integer",
			}
			break
		}

		tagID := uint(parsed)
		if _, ok := seen[tagID]; ok {
			continue
		}
		seen[tagID] = struct{}{}
		tagIDs = append(tagIDs, tagID)
	}

	if len(details) > 0 {
		return nil, apperror.WithDetails(apperror.ErrValidation, "invalid query parameters", details)
	}

	return tagIDs, nil
}
