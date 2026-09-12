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

	"ngumpul-host/backend/internal/access"
	"ngumpul-host/backend/internal/activity"
	"ngumpul-host/backend/internal/admin"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/availability"
	"ngumpul-host/backend/internal/comment"
	"ngumpul-host/backend/internal/config"
	"ngumpul-host/backend/internal/database"
	"ngumpul-host/backend/internal/hosting"
	"ngumpul-host/backend/internal/notification"
	"ngumpul-host/backend/internal/project"
	"ngumpul-host/backend/internal/report"
	"ngumpul-host/backend/internal/storage"
	"ngumpul-host/backend/internal/system"
	"ngumpul-host/backend/internal/user"
	"ngumpul-host/backend/internal/visits"
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
	accessService := access.NewService(pool, cfg)
	accessHandler := access.NewHandler(accessService, cfg)

	sessionManager := auth.NewSessionManager(pool)
	authHandler := auth.NewHandler(pool, sessionManager, cfg, accessService)

	storageService, err := storage.NewService(cfg, pool)
	if err != nil {
		log.Fatalf("Fatal: Storage initialization failed: %v", err)
	}

	availService := availability.NewService(pool)
	hostSpecs := system.GetHostSpecs()
	if err := availService.Initialize(ctx, int64(hostSpecs.UptimeSeconds)); err != nil {
		log.Printf("[Availability] Initialization warning: %v", err)
	}
	availService.StartHeartbeatWorker(ctx, 5*time.Minute)
	availService.StartProjectCheckWorker(ctx, 5*time.Minute)

	visitsHandler := visits.NewHandler(pool)
	commentHandler := comment.NewHandler(pool)
	reportHandler := report.NewHandler(pool)

	systemHandler := system.NewHandler(pool, availService)
	userHandler := user.NewHandler(pool)
	projectHandler := project.NewHandler(pool, availService, visitsHandler)
	hostingHandler := hosting.NewHandler(pool)
	activityHandler := activity.NewHandler(pool)
	notificationHandler := notification.NewHandler(pool)
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
	r.Use(auth.Middleware(sessionManager))

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

	// Outbound Visit Redirect ("Visits from Ngumpul" outbound click counting)
	r.Get("/go/{slug}", visitsHandler.Redirect)

	// Static uploaded files
	r.Handle("/uploads/*", storageService.ServeHandler())

	// Core API Router (mounted at /api and /api/v1 for full backwards-compatibility)
	apiRouter := chi.NewRouter()
	apiRouter.Get("/health", healthHandler)

	// Public Auth & Access
	apiRouter.Get("/auth/mode", accessHandler.GetRegistrationMode)
	apiRouter.Get("/invitations/validate", accessHandler.ValidateInvitation)
	apiRouter.Post("/auth/register", authHandler.Register)
	apiRouter.Post("/auth/login", authHandler.Login)
	apiRouter.Post("/auth/logout", authHandler.Logout)
	apiRouter.Get("/me", authHandler.Me)

	// Public Resources & Telemetry
	apiRouter.Get("/projects", projectHandler.ListPublic)
	apiRouter.Get("/projects/{slug}", projectHandler.GetBySlug)
	apiRouter.Get("/projects/{slug}/comments", commentHandler.ListByProject)
	apiRouter.Get("/users", userHandler.ListMembers)
	apiRouter.Get("/users/{username}", userHandler.GetMember)
	apiRouter.Get("/activity", activityHandler.ListPublic)
	apiRouter.Get("/status", systemHandler.GetPublicStatus)
	apiRouter.Get("/public/server", systemHandler.GetPublicServer)

	// Authenticated User Routes
	apiRouter.Group(func(userRouter chi.Router) {
		userRouter.Use(auth.RequireAuth)

		// Personal space & profile
		userRouter.Patch("/me", authHandler.UpdateProfile)
		userRouter.Get("/me/projects", projectHandler.ListMyProjects)
		userRouter.Get("/me/projects/{id}", projectHandler.GetMyProject)
		userRouter.Get("/me/projects/{id}/visits", visitsHandler.GetProjectVisits)
		userRouter.Patch("/me/projects/{id}", projectHandler.UpdateMyProject)

		// Comments & Reporting
		userRouter.Post("/projects/{slug}/comments", commentHandler.Create)
		userRouter.Delete("/comments/{id}", commentHandler.Delete)
		userRouter.Post("/comments/{id}/report", reportHandler.ReportComment)
		userRouter.Post("/projects/{slug}/report", reportHandler.ReportProject)

		// Hosting requests
		userRouter.Post("/hosting-requests", hostingHandler.Submit)
		userRouter.Get("/me/hosting-requests", hostingHandler.ListMyRequests)

		// Notifications
		userRouter.Get("/me/notifications", notificationHandler.ListMyNotifications)
		userRouter.Patch("/me/notifications/{id}/read", notificationHandler.MarkRead)
		userRouter.Post("/me/notifications/read-all", notificationHandler.MarkAllRead)

		// Activity
		userRouter.Get("/me/activity", activityHandler.ListMyActivity)

		// File upload (avatars, covers)
		userRouter.Post("/upload", storageService.UploadHandler)
	})

	// Administrator Routes
	apiRouter.Group(func(adminRouter chi.Router) {
		adminRouter.Use(auth.RequireAdmin)

		adminRouter.Get("/admin/stats", adminHandler.GetStats)

		// Member management
		adminRouter.Get("/admin/users", adminHandler.ListUsers)
		adminRouter.Patch("/admin/users/{id}/role", adminHandler.UpdateUserRole)
		adminRouter.Patch("/admin/users/{id}/status", adminHandler.UpdateUserStatus)

		// Access control & invitations
		adminRouter.Get("/admin/settings", accessHandler.GetSettings)
		adminRouter.Patch("/admin/settings", accessHandler.UpdateSettings)
		adminRouter.Get("/admin/invitations", accessHandler.ListInvitations)
		adminRouter.Post("/admin/invitations", accessHandler.CreateInvitation)
		adminRouter.Post("/admin/invitations/{id}/revoke", accessHandler.RevokeInvitation)

		// Project administration
		adminRouter.Get("/admin/projects", projectHandler.AdminList)
		adminRouter.Post("/admin/projects", projectHandler.AdminCreate)
		adminRouter.Patch("/admin/projects/{id}", projectHandler.AdminUpdate)

		// Moderation: Comments & Reports
		adminRouter.Get("/admin/comments", commentHandler.AdminList)
		adminRouter.Delete("/admin/comments/{id}", commentHandler.Delete)
		adminRouter.Get("/admin/reports", reportHandler.AdminList)
		adminRouter.Post("/admin/reports/{id}/resolve", reportHandler.AdminResolve)
		adminRouter.Post("/admin/reports/{id}/dismiss", reportHandler.AdminDismiss)

		// Hosting requests review
		adminRouter.Get("/admin/hosting-requests", hostingHandler.AdminListRequests)
		adminRouter.Post("/admin/hosting-requests/{id}/approve", hostingHandler.AdminApprove)
		adminRouter.Post("/admin/hosting-requests/{id}/reject", hostingHandler.AdminReject)
		adminRouter.Post("/admin/hosting-requests/{id}/complete", hostingHandler.AdminComplete)

		// Administrative feeds & system telemetry
		adminRouter.Get("/admin/activity", activityHandler.ListAdmin)
		adminRouter.Get("/admin/system", adminHandler.GetSystemHealth)
		adminRouter.Get("/admin/audit", adminHandler.ListAuditLogs)
	})

	// Mount API routes at /api (canonical) and /api/v1 (compat)
	r.Mount("/api", apiRouter)
	r.Mount("/api/v1", apiRouter)

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

	availService.Stop()
	log.Println("Backend server exited cleanly.")
}
