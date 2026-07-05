package validation

import (
	"regexp"
	"strings"

	"backend/internal/apperror"
)

var hexColorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

func NormalizeTagName(name string) string {
	return strings.TrimSpace(name)
}

func NormalizeTagColor(color *string) *string {
	if color == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*color)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func ValidateTag(name string, color *string) error {
	details := make(map[string]any)

	if name == "" {
		details["name"] = "name is required"
	} else if len([]rune(name)) > 40 {
		details["name"] = "name must be at most 40 characters"
	}

	if color != nil && !hexColorPattern.MatchString(*color) {
		details["color"] = "color must match #RRGGBB"
	}

	if len(details) > 0 {
		return apperror.WithDetails(apperror.ErrValidation, "invalid tag input", details)
	}

	return nil
}
