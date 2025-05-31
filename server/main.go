package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gorhandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"serve/config"
	"serve/database"
	"serve/handlers"
	"serve/middleware"
	"serve/services"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading env file: ", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection and migrate
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize email and text services
	emailService := services.NewEmailService(cfg)
	textService := services.NewTextService(cfg)

	// Initialize maps service
	mapsService := services.NewMapsService()

	// Initialize scheduler service
	scheduler := services.NewScheduler(db, emailService, textService)
	go scheduler.Start()
	defer scheduler.Stop()

	// Create a new router
	r := mux.NewRouter()

	// Initialize rate limiter (60 requests per minute per IP)
	rateLimiter := middleware.NewIPRateLimiter(rate.Every(time.Minute/60), 60)

	// Start cleanup routine for rate limiter
	go rateLimiter.CleanupOldEntries()

	// Set up middleware
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.RateLimitMiddleware(rateLimiter))

	// API routes (with auth)
	api := r.PathPrefix("/api").Subrouter()

	// User routes
	userRouter := api.PathPrefix("/users").Subrouter()
	handlers.RegisterUserRoutes(userRouter, db, emailService)

	// Project routes
	projectRouter := api.PathPrefix("/projects").Subrouter()
	handlers.RegisterProjectRoutes(projectRouter, db, cfg, emailService, textService)

	// Admin routes
	adminRouter := api.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware(cfg))
	adminRouter.Use(middleware.AdminMiddleware)
	handlers.RegisterAdminRoutes(adminRouter, db)

	// Geocoding routes
	geocodingHandler := &handlers.GeocodingHandler{
		MapsService: mapsService,
	}
	api.HandleFunc("/geocode", geocodingHandler.GeocodeAddress).Methods("POST")

	origin := "http://localhost:" + cfg.ServerPort
	corsHandler := gorhandler.CORS(
		gorhandler.AllowedOrigins(
			[]string{"http://localhost:3000", "http://localhost:5000", origin},
		), // Allowed origins
		gorhandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),           // Allowed methods
		gorhandler.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), // Allowed headers
		gorhandler.AllowCredentials(), // Allow credentials
	)(r)

	// Server setup
	address := ":" + cfg.ServerPort
	srv := &http.Server{
		Handler:      corsHandler,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server started on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
