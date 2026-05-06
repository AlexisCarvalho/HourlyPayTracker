package main

import (
	"embed"
	"io/fs"
	"net/http"
	"time"

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

	config.DB.AutoMigrate(&model.CompanyPaymentInformation{})
	config.DB.AutoMigrate(&model.User{})
	config.DB.AutoMigrate(&model.TimeEntry{})

	// =========================
	// GIN
	// =========================
	r := gin.Default()

	// CORS PRIMEIRO e MÁXIMO PERMISSIVO (antes de qualquer rota ou middleware)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = false // false com AllowAllOrigins
	corsConfig.MaxAge = 12 * time.Hour
	r.Use(cors.New(corsConfig))

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

	companyPaymentInfoRepo := repository.NewCompanyPaymentInformationRepository(config.DB)
	companyPaymentInfoService := service.NewCompanyPaymentInformationService(companyPaymentInfoRepo)
	companyPaymentInfoHandler := handler.NewCompanyPaymentInformationHandler(companyPaymentInfoService)

	timeEntryRepo := repository.NewTimeEntryRepository(config.DB)
	timeEntryService := service.NewTimeEntryService(timeEntryRepo)
	timeEntryHandler := handler.NewTimeEntryHandler(timeEntryService, userService)

	// =========================
	// API ENDPOINTS
	// =========================

	userHandler.RegisterRoutes(r)
	companyPaymentInfoHandler.RegisterRoutes(r)
	timeEntryHandler.RegisterRoutes(r)

	// =========================
	// START SERVER
	// =========================
	r.Run(":8080")
}
