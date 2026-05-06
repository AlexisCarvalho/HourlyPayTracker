package handler

import (
	"calcula_pagamento/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CompanyPaymentInformationHandler struct {
	service service.CompanyPaymentInformationService
}

func NewCompanyPaymentInformationHandler(s service.CompanyPaymentInformationService) *CompanyPaymentInformationHandler {
	return &CompanyPaymentInformationHandler{service: s}
}

func (h *CompanyPaymentInformationHandler) RegisterRoutes(r *gin.Engine) {
	// Public endpoint - can list companies without auth
	r.GET("/companies", h.GetAll)

	// For future: authenticated endpoints
	// api := r.Group("/companies")
	// api.Use(middleware.AuthMiddleware())
}

func (h *CompanyPaymentInformationHandler) GetAll(c *gin.Context) {
	companies, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, companies)
}
