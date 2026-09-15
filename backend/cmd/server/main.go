package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"opsdesk/internal/config"
	"opsdesk/internal/database"
	"opsdesk/internal/handlers"
	"opsdesk/internal/middleware"
	"opsdesk/internal/models"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting OpsDesk Backend Service on port %s...", cfg.Port)

	db, err := database.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Fatal: Database initialization failed: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()

	// Global Middlewares
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.LoggerMiddleware)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORSMiddleware(cfg))

	// Rate limiter (60 requests per minute)
	rateLimiter := middleware.NewRateLimiter(60, time.Minute)
	r.Use(rateLimiter.Limit)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(cfg)
	ticketHandler := handlers.NewTicketHandler()
	commentHandler := handlers.NewCommentHandler()
	assetHandler := handlers.NewAssetHandler()
	attachmentHandler := handlers.NewAttachmentHandler(cfg)
	adminHandler := handlers.NewAdminHandler()
	systemHandler := handlers.NewSystemHandler(cfg)
	diagnosticsHandler := handlers.NewDiagnosticsHandler()
	webhookHandler := handlers.NewWebhookHandler()

	// Public Routes
	r.Get("/api/health", systemHandler.HealthCheck)
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)
	r.Post("/api/auth/logout", authHandler.Logout)

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg))

		// User Profile
		r.Get("/api/auth/me", authHandler.Me)
		r.Put("/api/auth/profile", authHandler.UpdateProfile)

		// Tickets API
		r.Get("/api/tickets", ticketHandler.List)
		r.Get("/api/tickets/search", ticketHandler.Search)
		r.Post("/api/tickets", ticketHandler.Create)
		r.Get("/api/tickets/{id}", ticketHandler.Get)
		r.Put("/api/tickets/{id}", ticketHandler.Update)
		r.Delete("/api/tickets/{id}", ticketHandler.Delete)

		// Ticket Comments
		r.Get("/api/tickets/{id}/comments", commentHandler.ListByTicket)
		r.Post("/api/tickets/{id}/comments", commentHandler.Create)
		r.Delete("/api/tickets/{id}/comments/{commentId}", commentHandler.Delete)

		// Ticket Attachments
		r.Get("/api/tickets/{id}/attachments", attachmentHandler.List)
		r.Post("/api/tickets/{id}/attachments", attachmentHandler.Upload)
		r.Post("/api/tickets/{id}/attachments/bulk", attachmentHandler.BulkUpload)
		r.Get("/api/tickets/{id}/attachments/{attachmentId}/download", attachmentHandler.Download)
		r.Get("/api/attachments/preview", attachmentHandler.Preview)

		// Assets API
		r.Get("/api/assets", assetHandler.List)
		r.Get("/api/assets/search", assetHandler.Search)
		r.Get("/api/assets/export", assetHandler.ExportReport)
		r.Post("/api/assets", assetHandler.Create)
		r.Post("/api/assets/import-xml", assetHandler.ImportXML)
		r.Post("/api/assets/import-json", assetHandler.ImportJSON)
		r.Get("/api/assets/{id}", assetHandler.Get)
		r.Put("/api/assets/{id}", assetHandler.Update)
		r.Delete("/api/assets/{id}", assetHandler.Delete)

		// Diagnostics & Network Tools
		r.Post("/api/assets/{id}/ping", diagnosticsHandler.PingAsset)
		r.Post("/api/diagnostics/ping", diagnosticsHandler.PingAsset)
		r.Post("/api/diagnostics/port-check", diagnosticsHandler.PortCheck)

		// Webhooks & Integrations
		r.Post("/api/webhooks/test", webhookHandler.TestWebhook)

		// Admin Routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole(models.RoleAdmin))
			r.Get("/api/admin/users", adminHandler.ListUsers)
			r.Put("/api/admin/users/{id}/role", adminHandler.UpdateUserRole)
			r.Get("/api/admin/audit-logs", adminHandler.ListAuditLogs)
			r.Get("/api/admin/system/metrics", systemHandler.Metrics)
		})
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("OpsDesk API Server listening on http://localhost:%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
