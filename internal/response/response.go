package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeValidationError = "VALIDATION_ERROR"
	CodeNotFound        = "NOT_FOUND"
	CodeConflict        = "CONFLICT"
	CodeDatabaseError   = "DATABASE_ERROR"
	CodeInternalError   = "INTERNAL_ERROR"
)

type Envelope struct {
	Data  any        `json:"data"`
	Error *ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{
		Data:  data,
		Error: nil,
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{
		Data:  data,
		Error: nil,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, status int, code string, message string, details map[string]any) {
	c.JSON(status, Envelope{
		Data: nil,
		Error: &ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func ValidationError(c *gin.Context, message string, details map[string]any) {
	Error(c, http.StatusBadRequest, CodeValidationError, message, details)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, CodeNotFound, message, nil)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, CodeConflict, message, nil)
}

func DatabaseError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, CodeDatabaseError, message, nil)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, CodeInternalError, message, nil)
}
