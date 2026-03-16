package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"calcula_pagamento/internal/middleware"
	"calcula_pagamento/internal/model"
	"calcula_pagamento/internal/service"

	"github.com/gin-gonic/gin"
)

type TimeEntryHandler struct {
	service service.TimeEntryService
}

type BulkMarkRequest struct {
	IDs []uint `json:"ids"`
}

func NewTimeEntryHandler(s service.TimeEntryService) *TimeEntryHandler {
	return &TimeEntryHandler{service: s}
}

func (h *TimeEntryHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/time-entry")
	api.Use(middleware.AuthMiddleware())

	// Creation and listing
	api.POST("/", h.CreateTimeEntry)
	api.GET("/durations", h.GetDurations)
	api.GET("/durations-month", h.GetDurationsMonth)
	api.GET("/pages/count", h.GetPagesCount)

	// Operations specific by ID
	api.PUT("/:id", h.Update)
	api.PUT("/:id/mark-paid", h.MarkAsPaid)

	// Batch operations (without ID)
	api.PUT("/mark-paid", h.MarkMultipleAsPaid)
	api.DELETE("/delete", h.Delete)
}

func (h *TimeEntryHandler) CreateTimeEntry(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	var entry model.TimeEntry

	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry.IdUser = userID

	if err := h.service.RegisterTime(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *TimeEntryHandler) MarkAsPaid(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.service.MarkAsPaid(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "marked as paid"})
}

func (h *TimeEntryHandler) Update(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	idParam := c.Param("id")
	idParsed, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	id := uint(idParsed)

	var entry model.TimeEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ensure the ID in URL and body (if provided) match
	if entry.ID != 0 && entry.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID in body does not match URL"})
		return
	}
	entry.ID = id
	entry.IdUser = userID

	if entry.ClockOut.Before(entry.ClockIn) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ClockOut must be after ClockIn"})
		return
	}

	if err := h.service.Update(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *TimeEntryHandler) MarkMultipleAsPaid(c *gin.Context) {
	var req BulkMarkRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or empty request"})
		return
	}

	if err := h.service.MarkMultipleAsPaid(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "entries marked as paid", "count": len(req.IDs)})
}

func (h *TimeEntryHandler) Delete(c *gin.Context) {
	var req BulkMarkRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or empty request"})
		return
	}

	if err := h.service.Delete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "entries deleted", "count": len(req.IDs)})
}

func (h *TimeEntryHandler) GetDurations(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idUser, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	paidParam := c.Query("paid")
	paid := false
	if paidParam != "" {
		p, err := strconv.ParseBool(paidParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid paid query parameter"})
			return
		}
		paid = p
	}

	descParam := c.Query("desc")
	desc := false
	if descParam != "" {
		d, err := strconv.ParseBool(descParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid desc query parameter"})
			return
		}
		desc = d
	}

	limit := 20
	page := 1
	if l := c.Query("limit"); l != "" {
		if lInt, err := strconv.Atoi(l); err == nil && lInt > 0 {
			limit = lInt
		} else if l == "-1" {
			limit = -1
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit query parameter"})
			return
		}
	}
	if p := c.Query("page"); p != "" {
		if pInt, err := strconv.Atoi(p); err == nil && pInt > 0 {
			page = pInt
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page query parameter"})
			return
		}
	}

	entries, err := h.service.GetDurations(idUser, paid, desc, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *TimeEntryHandler) GetDurationsMonth(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idUser, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	var year int
	var month int

	if y := c.Query("year"); y != "" {
		if yInt, err := strconv.Atoi(y); err == nil && yInt > 0 {
			year = yInt
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year query parameter"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you need to inform a year to query"})
		return
	}

	if m := c.Query("month"); m != "" {
		if mInt, err := strconv.Atoi(m); err == nil && mInt > 0 {
			month = mInt
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month query parameter"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you need to inform a month to query"})
		return
	}

	entries, err := h.service.GetDurationsMonth(idUser, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *TimeEntryHandler) GetPagesCount(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	idUser, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	paidParam := c.Query("paid")
	paid := false
	if paidParam != "" {
		p, err := strconv.ParseBool(paidParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid paid query parameter"})
			return
		}
		paid = p
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if lInt, err := strconv.Atoi(l); err == nil && lInt > 0 {
			limit = lInt
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit query parameter"})
			return
		}
	}

	totalEntries, err := h.service.CountEntries(idUser, paid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	totalPages := (totalEntries + limit - 1) / limit // ceil
	c.JSON(http.StatusOK, gin.H{"total_pages": totalPages, "total_entries": totalEntries, "limit": limit})
}
