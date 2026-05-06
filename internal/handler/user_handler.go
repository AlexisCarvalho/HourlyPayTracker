package handler

import (
	"calcula_pagamento/internal/auth"
	"calcula_pagamento/internal/middleware"
	"calcula_pagamento/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"

	service "calcula_pagamento/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/users/login", h.Login)
	r.POST("/users/register", h.RegisterUser)

	api := r.Group("/users")
	api.Use(middleware.AuthMiddleware())
	api.GET("/me", h.GetMe)
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var entry model.User

	println(entry.PreferredCompanyID)
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RegisterUser(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Code     string `json:"code"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	user := h.service.Authenticate(req.Code, req.Password)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais inválidas"})
		return
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	user, err := h.service.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a response without the password
	response := gin.H{
		"id":                   user.ID,
		"name":                 user.Name,
		"code":                 user.Code,
		"preferred_company_id": user.PreferredCompanyID,
		"preferred_company":    user.PreferredCompany,
	}

	c.JSON(http.StatusOK, response)
}
