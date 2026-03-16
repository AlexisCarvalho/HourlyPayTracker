package main

import (
	"embed"
	"io/fs"
	"net/http"

	"calcula_pagamento/internal/config"
	"calcula_pagamento/internal/handler"
	"calcula_pagamento/internal/model"
	"calcula_pagamento/internal/repository"
	service "calcula_pagamento/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//
// EMBED OF FRONTEND
//

//go:embed static/*
var staticFiles embed.FS

func main() {
	// =========================
	// INITIAL CONFIGURATION
	// =========================
	config.LoadEnv()
	config.ConnectDB()

	config.DB.AutoMigrate(&model.User{})
	config.DB.AutoMigrate(&model.TimeEntry{})

	// =========================
	// GIN
	// =========================
	r := gin.Default()
	r.Use(cors.Default())

	// =========================
	// FRONTEND (SPA)
	// =========================

	// Remove the prefix "static" of embed
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic("erro ao carregar arquivos estáticos")
	}

	// Static Files (JS, CSS, assets...)
	r.StaticFS("/account", http.FS(subFS))
	r.StaticFS("/dashboard", http.FS(subFS))

	// Serve account.html directly in /account
	r.GET("/account", func(c *gin.Context) {
		c.FileFromFS("account.html", http.FS(subFS))
	})

	// Serve dashboard.html directly in /dashboard
	r.GET("/dashboard", func(c *gin.Context) {
		c.FileFromFS("dashboard.html", http.FS(subFS))
	})

	// Redirect / to /account
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/account")
	})

	// SPA fallback
	r.NoRoute(func(c *gin.Context) {
		c.FileFromFS("account.html", http.FS(subFS))
	})

	// =========================
	// REPOSITORIES / SERVICES
	// =========================
	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	timeEntryRepo := repository.NewTimeEntryRepository(config.DB)
	timeEntryService := service.NewTimeEntryService(timeEntryRepo)
	timeEntryHandler := handler.NewTimeEntryHandler(timeEntryService)

	// =========================
	// API ENDPOINTS
	// =========================

	userHandler.RegisterRoutes(r)
	timeEntryHandler.RegisterRoutes(r)

	// =========================
	// START SERVER
	// =========================
	r.Run(":8080")
}
