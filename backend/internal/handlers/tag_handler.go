package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"project-alpha/backend/internal/apperror"
	"project-alpha/backend/internal/dto"
	"project-alpha/backend/internal/response"
	"project-alpha/backend/internal/services"
)

type TagHandler struct {
	tags *services.TagService
}

func NewTagHandler(tags *services.TagService) *TagHandler {
	return &TagHandler{tags: tags}
}

func (h *TagHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/tags", h.Create)
	group.GET("/tags", h.List)
	group.PUT("/tags/:id", h.Update)
	group.DELETE("/tags/:id", h.Delete)
}

func (h *TagHandler) Create(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "invalid request body", map[string]any{"body": err.Error()})
		return
	}

	tag, err := h.tags.Create(req)
	if err != nil {
		respondError(c, err)
		return
	}

	response.Created(c, dto.NewTagResponse(tag))
}

func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.tags.List()
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, dto.NewTagResponses(tags))
}

func (h *TagHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "invalid request body", map[string]any{"body": err.Error()})
		return
	}

	tag, err := h.tags.Update(id, req)
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, dto.NewTagResponse(tag))
}

func (h *TagHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.tags.Delete(id); err != nil {
		respondError(c, err)
		return
	}

	response.NoContent(c)
}

func parseIDParam(c *gin.Context, name string) (uint, bool) {
	raw := c.Param(name)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		respondError(c, apperror.WithDetails(
			apperror.ErrValidation,
			"invalid path parameter",
			map[string]any{name: "must be a positive integer"},
		))
		return 0, false
	}

	return uint(id), true
}
