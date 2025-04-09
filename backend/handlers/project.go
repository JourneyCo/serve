package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"project-registration-system/middleware"
	"project-registration-system/models"
	"project-registration-system/services"
)

// ProjectHandler handles project-related requests
type ProjectHandler struct {
	DB           *sql.DB
	EmailService *services.EmailService
}

// regRequest defines the JSON request for registration
type regRequest struct {
	GuestCount       int    `json:"guest_count"`
	IsLeadInterested bool   `json:"lead_interest"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
}

// RegisterProjectRoutes registers the routes for project handlers
func RegisterProjectRoutes(router *mux.Router, db *sql.DB, emailService *services.EmailService) {
	handler := &ProjectHandler{
		DB:           db,
		EmailService: emailService,
	}

	router.HandleFunc("", handler.GetProjects).Methods("GET")
	router.HandleFunc("/{id:[0-9]+}", handler.GetProject).Methods("GET")
	router.HandleFunc("/{id:[0-9]+}/register", handler.RegisterForProject).Methods("POST")
	router.HandleFunc("/{id:[0-9]+}/cancel", handler.CancelRegistration).Methods("POST")
	router.HandleFunc("/{id:[0-9]+}/registrations", handler.GetProjectRegistrations).Methods("GET")
}

// GetProjects returns all projects
func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := models.GetAllProjects(h.DB)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve projects")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, projects)
}

// GetProject returns a specific project by ID
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := models.GetProjectByID(h.DB, id)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve project")
		return
	}

	if project == nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	project.Tools = []models.Tool{}

	middleware.RespondWithJSON(w, http.StatusOK, project)
}

// RegisterForProject registers a user for a project
func (h *ProjectHandler) RegisterForProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Get user ID from the token
	userID, err := middleware.GetUserIDFromRequest(r)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Failed to get user information")
		return
	}

	// Get or create user in database
	user, err := models.GetUserByID(h.DB, userID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to check user status")
		return
	}

	// If user doesn't exist, create them
	if user == nil {
		user = &models.User{
			ID: userID,
		}
		err = models.CreateUser(h.DB, user)
		if err != nil {
			middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}
	}

	var reg regRequest
	// Parse registration request
	if err = middleware.ParseJSON(r, &reg); err != nil {
		// If parsing fails, use default values (for backward compatibility)
		reg.GuestCount = 0
		reg.IsLeadInterested = false
	}
	err = nil

	// Update user information
	user.FirstName = reg.FirstName
	user.LastName = reg.LastName
	user.Phone = reg.Phone
	user.Email = reg.Email
	if err = models.UpdateUser(h.DB, user); err != nil {
		log.Printf("Failed to update user information: %v\n", err)
		middleware.RespondWithError(
			w, http.StatusInternalServerError, "Failed to update user information",
		)
		return
	}

	// Validate guest count
	if reg.GuestCount < 0 {
		middleware.RespondWithError(w, http.StatusBadRequest, "Guest count cannot be negative")
		return
	}

	// Register for the project
	registration, err := models.RegisterForProject(
		h.DB, userID, projectID, reg.GuestCount, reg.IsLeadInterested,
	)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get project details for email
	project, err := models.GetProjectByID(h.DB, projectID)
	if err != nil {
		middleware.RespondWithError(
			w, http.StatusInternalServerError, "Registration successful but failed to send confirmation email",
		)
		return
	}

	// Get user details for email
	user, err = models.GetUserByID(h.DB, userID)
	if err != nil {
		middleware.RespondWithError(
			w, http.StatusInternalServerError, "Registration successful but failed to send confirmation email",
		)
		return
	}

	// Send confirmation email
	if user != nil && project != nil {
		go h.EmailService.SendRegistrationConfirmation(user, project)
	}

	middleware.RespondWithJSON(w, http.StatusCreated, registration)
}

// CancelRegistration cancels a user's registration for a project
func (h *ProjectHandler) CancelRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Get user ID from the token
	userID, err := middleware.GetUserIDFromRequest(r)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Failed to get user information")
		return
	}

	// Cancel the registration
	err = models.CancelRegistration(h.DB, userID, projectID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Registration cancelled successfully"})
}

// GetProjectRegistrations returns all registrations for a project (admin only)
func (h *ProjectHandler) GetProjectRegistrations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	registrations, err := models.GetProjectRegistrations(h.DB, projectID)
	if err != nil {
		log.Println("failed to get project registrations")
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve registrations")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, registrations)
}
