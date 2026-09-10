package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ngumpul-host/backend/internal/activities"
	"ngumpul-host/backend/internal/admin"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/config"
	"ngumpul-host/backend/internal/database"
	"ngumpul-host/backend/internal/notifications"
	"ngumpul-host/backend/internal/projects"
	"ngumpul-host/backend/internal/requests"
	"ngumpul-host/backend/internal/storage"
	"ngumpul-host/backend/internal/users"
)

func main() {
	log.Println("Starting Ngumpul Host Backend Service...")

	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database and apply migrations
	pool, err := database.ConnectAndMigrate(ctx, cfg)
	if err != nil {
		log.Fatalf("Fatal: Database initialization failed: %v", err)
	}
	defer pool.Close()

	// Initialize services
	sessionManager := auth.NewSessionManager(pool)
	authHandler := auth.NewHandler(pool, sessionManager, cfg)

	storageService, err := storage.NewStorageService(cfg)
	if err != nil {
		log.Fatalf("Fatal: Storage initialization failed: %v", err)
	}

	usersHandler := users.NewHandler(pool)
	projectsHandler := projects.NewHandler(pool)
	requestsHandler := requests.NewHandler(pool)
	activitiesHandler := activities.NewHandler(pool)
	notificationsHandler := notifications.NewHandler(pool)
	adminHandler := admin.NewHandler(pool)

	r := chi.NewRouter()

	// Standard middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration for local development / reverse proxy
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost*", "http://127.0.0.1*", cfg.AppURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Session extraction middleware
	r.Use(auth.AuthMiddleware(sessionManager))

	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "ngumpul-host-backend",
			"time":    time.Now().Format(time.RFC3339),
		})
	}

	// Health check endpoint
	r.Get("/health", healthHandler)

	// Static uploaded files
	r.Handle("/uploads/*", storageService.ServeHandler())

	// API v1 Routes
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/health", healthHandler)
		// Public Auth
		api.Post("/auth/register", authHandler.Register)
		api.Post("/auth/login", authHandler.Login)
		api.Post("/auth/logout", authHandler.Logout)
		api.Get("/me", authHandler.Me)

		// Public Resources
		api.Get("/projects", projectsHandler.ListPublic)
		api.Get("/projects/{slug}", projectsHandler.GetBySlug)
		api.Get("/users", usersHandler.ListMembers)
		api.Get("/users/{username}", usersHandler.GetMember)
		api.Get("/activity", activitiesHandler.ListPublic)
		api.Get("/status", adminHandler.GetPublicStatus)

		// Authenticated User Routes
		api.Group(func(userRouter chi.Router) {
			userRouter.Use(auth.RequireAuth)

			// Personal space & profile
			userRouter.Patch("/me", authHandler.UpdateProfile)
			userRouter.Get("/me/projects", projectsHandler.ListMyProjects)
			userRouter.Patch("/me/projects/{id}", projectsHandler.UpdateMyProject)

			// Hosting requests
			userRouter.Post("/hosting-requests", requestsHandler.Submit)
			userRouter.Get("/me/hosting-requests", requestsHandler.ListMyRequests)

			// Notifications
			userRouter.Get("/me/notifications", notificationsHandler.ListMyNotifications)
			userRouter.Patch("/me/notifications/{id}/read", notificationsHandler.MarkRead)
			userRouter.Post("/me/notifications/read-all", notificationsHandler.MarkAllRead)

			// File upload (avatars, covers)
			userRouter.Post("/upload", storageService.UploadHandler)
		})

		// Administrator Routes
		api.Group(func(adminRouter chi.Router) {
			adminRouter.Use(auth.RequireAdmin)

			adminRouter.Get("/admin/stats", adminHandler.GetStats)

			// Member management
			adminRouter.Get("/admin/users", adminHandler.ListUsers)
			adminRouter.Patch("/admin/users/{id}/role", adminHandler.UpdateUserRole)
			adminRouter.Patch("/admin/users/{id}/status", adminHandler.UpdateUserStatus)

			// Project administration
			adminRouter.Get("/admin/projects", projectsHandler.AdminList)
			adminRouter.Post("/admin/projects", projectsHandler.AdminCreate)
			adminRouter.Patch("/admin/projects/{id}", projectsHandler.AdminUpdate)

			// Hosting requests review
			adminRouter.Get("/admin/hosting-requests", requestsHandler.AdminListRequests)
			adminRouter.Post("/admin/hosting-requests/{id}/approve", requestsHandler.AdminApprove)
			adminRouter.Post("/admin/hosting-requests/{id}/reject", requestsHandler.AdminReject)
			adminRouter.Post("/admin/hosting-requests/{id}/complete", requestsHandler.AdminComplete)

			// Administrative feeds & system telemetry
			adminRouter.Get("/admin/activity", activitiesHandler.ListAdmin)
			adminRouter.Get("/admin/system", adminHandler.GetSystemHealth)
			adminRouter.Get("/admin/audit", adminHandler.ListAuditLogs)
		})
	})

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Ngumpul Host Backend listening on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server ListenAndServe error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down backend gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Backend server exited cleanly.")
}
