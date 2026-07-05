package handlers

import (
	"github.com/gin-gonic/gin"

	"backend/internal/dto"
	"backend/internal/response"
	"backend/internal/services"
)

type TicketHandler struct {
	tickets *services.TicketService
}

func NewTicketHandler(tickets *services.TicketService) *TicketHandler {
	return &TicketHandler{tickets: tickets}
}

func (h *TicketHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/tickets", h.Create)
	group.GET("/tickets", h.List)
	group.GET("/tickets/:id", h.Get)
	group.PUT("/tickets/:id", h.Update)
	group.DELETE("/tickets/:id", h.Delete)
	group.POST("/tickets/:id/complete", h.Complete)
	group.POST("/tickets/:id/uncomplete", h.Uncomplete)
	group.POST("/tickets/:id/tags/:tagId", h.AddTag)
	group.DELETE("/tickets/:id/tags/:tagId", h.RemoveTag)
}

func (h *TicketHandler) Create(c *gin.Context) {
	var req dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "invalid request body", map[string]any{"body": err.Error()})
		return
	}

	ticket, err := h.tickets.Create(req)
	if err != nil {
		respondError(c, err)
		return
	}

	response.Created(c, dto.NewTicketResponse(ticket))
}

func (h *TicketHandler) List(c *gin.Context) {
	var query dto.ListTicketsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, "invalid query parameters", map[string]any{"query": err.Error()})
		return
	}

	tickets, pagination, err := h.tickets.List(query)
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, dto.TicketListResponse{
		Items:      dto.NewTicketResponses(tickets),
		Pagination: pagination,
	})
}

func (h *TicketHandler) Get(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	ticket, err := h.tickets.Get(id)
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, dto.NewTicketResponse(ticket))
}

func (h *TicketHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "invalid request body", map[string]any{"body": err.Error()})
		return
	}

	ticket, err := h.tickets.Update(id, req)
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, dto.NewTicketResponse(ticket))
}

func (h *TicketHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.tickets.Delete(id); err != nil {
		respondError(c, err)
		return
	}

	response.NoContent(c)
}

func (h *TicketHandler) Complete(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	ticket, err := h.tickets.Complete(id)
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, dto.NewTicketResponse(ticket))
}

func (h *TicketHandler) Uncomplete(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	ticket, err := h.tickets.Uncomplete(id)
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, dto.NewTicketResponse(ticket))
}

func (h *TicketHandler) AddTag(c *gin.Context) {
	ticketID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	tagID, ok := parseIDParam(c, "tagId")
	if !ok {
		return
	}

	ticket, err := h.tickets.AddTag(ticketID, tagID)
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, dto.NewTicketResponse(ticket))
}

func (h *TicketHandler) RemoveTag(c *gin.Context) {
	ticketID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	tagID, ok := parseIDParam(c, "tagId")
	if !ok {
		return
	}

	ticket, err := h.tickets.RemoveTag(ticketID, tagID)
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, dto.NewTicketResponse(ticket))
}
