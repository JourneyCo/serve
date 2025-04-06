package handlers

import (
	"database/sql"
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

	middleware.RespondWithJSON(w, http.StatusOK, project)
}

// RegistrationRequest defines the JSON request for registration
type RegistrationRequest struct {
	GuestCount    int  `json:"guestCount"`
	IsProjectLead bool `json:"isProjectLead"`
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

	// Parse registration request
	var regRequest struct {
		GuestCount    int    `json:"guestCount"`
		IsProjectLead bool   `json:"isProjectLead"`
		Phone         string `json:"phone"`
		ContactEmail  string `json:"contactEmail"`
	}
	if err := middleware.ParseJSON(r, &regRequest); err != nil {
		// If parsing fails, use default values (for backward compatibility)
		regRequest.GuestCount = 0
		regRequest.IsProjectLead = false
	}

	// Update user information
	user.Phone = regRequest.Phone
	user.ContactEmail = regRequest.ContactEmail
	err = models.UpdateUser(h.DB, user)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update user information")
		return
	}

	// Validate guest count
	if regRequest.GuestCount < 0 {
		middleware.RespondWithError(w, http.StatusBadRequest, "Guest count cannot be negative")
		return
	}

	// Register for the project
	registration, err := models.RegisterForProject(h.DB, userID, projectID, regRequest.GuestCount, regRequest.IsProjectLead)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get project details for email
	project, err := models.GetProjectByID(h.DB, projectID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Registration successful but failed to send confirmation email")
		return
	}

	// Get user details for email
	user, err = models.GetUserByID(h.DB, userID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Registration successful but failed to send confirmation email")
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
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve registrations")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, registrations)
}
