package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"project-alpha/backend/internal/apperror"
	"project-alpha/backend/internal/response"
)

func respondError(c *gin.Context, err error) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		switch {
		case apperror.Is(err, apperror.ErrValidation):
			response.ValidationError(c, appErr.Message, appErr.Details)
		case apperror.Is(err, apperror.ErrNotFound):
			response.NotFound(c, appErr.Message)
		case apperror.Is(err, apperror.ErrConflict):
			response.Conflict(c, appErr.Message)
		case apperror.Is(err, apperror.ErrDatabase):
			response.DatabaseError(c, appErr.Message)
		default:
			response.InternalError(c, appErr.Message)
		}
		return
	}

	response.InternalError(c, "internal server error")
}
