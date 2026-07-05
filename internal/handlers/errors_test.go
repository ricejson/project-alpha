package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"backend/internal/apperror"
)

func TestRespondErrorMapsApplicationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name   string
		err    error
		status int
		code   string
	}{
		{name: "validation", err: apperror.New(apperror.ErrValidation, "bad input"), status: http.StatusBadRequest, code: "VALIDATION_ERROR"},
		{name: "not found", err: apperror.New(apperror.ErrNotFound, "missing"), status: http.StatusNotFound, code: "NOT_FOUND"},
		{name: "conflict", err: apperror.New(apperror.ErrConflict, "duplicate"), status: http.StatusConflict, code: "CONFLICT"},
		{name: "database", err: apperror.New(apperror.ErrDatabase, "db failed"), status: http.StatusInternalServerError, code: "DATABASE_ERROR"},
		{name: "unknown", err: errors.New("boom"), status: http.StatusInternalServerError, code: "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/", func(c *gin.Context) {
				respondError(c, tt.err)
			})

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

			if recorder.Code != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, recorder.Code)
			}

			var body struct {
				Error *struct {
					Code string `json:"code"`
				} `json:"error"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if body.Error == nil || body.Error.Code != tt.code {
				t.Fatalf("expected code %s, got %+v", tt.code, body.Error)
			}
		})
	}
}

func TestParseIDParamRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/items/:id", func(c *gin.Context) {
		if _, ok := parseIDParam(c, "id"); ok {
			t.Fatal("expected invalid id")
		}
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/items/abc", nil))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", recorder.Code)
	}
}
