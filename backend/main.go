package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gorhandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"project-registration-system/config"
	"project-registration-system/database"
	"project-registration-system/handlers"
	"project-registration-system/middleware"
	"project-registration-system/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Println("Repl Name (Slug):", os.Getenv("REPL_SLUG"))
	fmt.Println("Repl Owner:", os.Getenv("REPL_OWNER"))
	fmt.Println("Repl ID:", os.Getenv("REPL_ID"))

	// Initialize database connection
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize email service
	emailService := services.NewEmailService(cfg)

	// Initialize maps service
	mapsService := services.NewMapsService()

	// Initialize scheduler service
	scheduler := services.NewScheduler(db, emailService)
	go scheduler.Start()
	defer scheduler.Stop()

	// Create a new router
	r := mux.NewRouter()

	// Set up middleware
	r.Use(middleware.LoggerMiddleware)

	// API routes (with auth)
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(cfg))

	// Public routes
	// No public routes needed for now
	r.PathPrefix("/api/public").Subrouter()

	// User routes
	userRouter := api.PathPrefix("/users").Subrouter()
	handlers.RegisterUserRoutes(userRouter, db, emailService)

	// Project routes
	projectRouter := api.PathPrefix("/projects").Subrouter()
	handlers.RegisterProjectRoutes(projectRouter, db, emailService)

	// Admin routes
	adminRouter := api.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AdminMiddleware)
	handlers.RegisterAdminRoutes(adminRouter, db)

	// Auth routes
	authRouter := r.PathPrefix("/auth").Subrouter()
	handlers.RegisterAuthRoutes(authRouter, cfg)

	// Geocoding routes
	geocodingHandler := &handlers.GeocodingHandler{
		MapsService: mapsService,
	}
	api.HandleFunc("/geocode", geocodingHandler.GeocodeAddress).Methods("POST")

	corsHandler := gorhandler.CORS(
		gorhandler.AllowedOrigins(
			[]string{"http://localhost:3000", "http://localhost:5000"},
		), // Allowed origins
		gorhandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),           // Allowed methods
		gorhandler.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), // Allowed headers
		gorhandler.AllowCredentials(), // Allow credentials
	)(r)

	// Server setup
	srv := &http.Server{
		Handler:      corsHandler,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server started on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
