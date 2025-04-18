package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"serve/middleware"
	"serve/models"
)

// AdminHandler handles admin-related requests
type AdminHandler struct {
	DB *sql.DB
}

// ProjectInput represents the input for creating or updating a project
type ProjectInput struct {
	Title                string  `json:"title"`
	Description          string  `json:"description"`
	ShortDescription     string  `json:"short_description"`
	Time                 string  `json:"time"`
	ProjectDate          string  `json:"project_date"`
	MaxCapacity          int     `json:"max_capacity"`
	WheelchairAccessible bool    `json:"wheelchair_accessible"`
	LeadUserID           string  `json:"lead_user_id"`
	Tools                []int   `json:"tools,omitempty"`
	Skills               []int   `json:"skills,omitempty"`
	Categories           []int   `json:"categories,omitempty"`
	Ages                 []int   `json:"ages,omitempty"`
	Supplies             []int   `json:"supplies,omitempty"`
	LocationName         string  `json:"location_name"`
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
}

// RegisterAdminRoutes registers the routes for admin handlers
func RegisterAdminRoutes(router *mux.Router, db *sql.DB) {
	handler := &AdminHandler{
		DB: db,
	}

	router.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	router.HandleFunc("/projects", handler.CreateProject).Methods("POST")
	router.HandleFunc("/projects/{id:[0-9]+}", handler.UpdateProject).Methods("PUT")
	router.HandleFunc("/projects/{id:[0-9]+}", handler.DeleteProject).Methods("DELETE")
}

// GetAllUsers returns all users
func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers(h.DB)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, users)
}

// CreateProject creates a new project
func (h *AdminHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var input ProjectInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println(err)
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if input.Title == "" || input.Description == "" || input.Time == "" || input.ProjectDate == "" || input.MaxCapacity <= 0 {
		middleware.RespondWithError(
			w, http.StatusBadRequest, "All fields are required and max capacity must be greater than 0",
		)
		return
	}

	// Parse project date
	projectDate, err := time.Parse("2006-01-02", input.ProjectDate)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project date format (use YYYY-MM-DD)")
		return
	}

	// Create project
	project := &models.Project{
		Title:            input.Title,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		Time:             input.Time,
		ProjectDate:      projectDate,
		MaxCapacity:      input.MaxCapacity,
		LocationName:     input.LocationName,
		Latitude:         input.Latitude,
		Longitude:        input.Longitude,
		LeadUserID:       input.LeadUserID,
	}

	project = applyAccessories(input, project)

	if err = models.CreateProject(h.DB, project); err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	middleware.RespondWithJSON(w, http.StatusCreated, project)
}

// UpdateProject updates an existing project
func (h *AdminHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Check if project exists
	project, err := models.GetProjectByID(h.DB, id)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve project")
		return
	}

	if project == nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	var input ProjectInput
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// Validate input
	if input.Title == "" || input.Description == "" || input.ShortDescription == "" || input.Time == "" || input.ProjectDate == "" || input.MaxCapacity <= 0 {
		middleware.RespondWithError(
			w, http.StatusBadRequest, "All fields are required and max capacity must be greater than 0",
		)
		return
	}

	// Parse project date
	if input.ProjectDate == "" {
		input.ProjectDate = "2025-07-12"
	}
	projectDate, err := time.Parse("2006-01-02", input.ProjectDate)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project date format (use YYYY-MM-DD)")
		return
	}

	// Update project
	project.Title = input.Title
	project.Description = input.Description
	project.ShortDescription = input.ShortDescription
	project.Time = input.Time
	project.ProjectDate = projectDate
	project.MaxCapacity = input.MaxCapacity
	project.LocationName = input.LocationName
	project.Latitude = input.Latitude
	project.Longitude = input.Longitude
	project.WheelchairAccessible = input.WheelchairAccessible

	// TODO: Need to delete accessories here for an existing projct before adding accessories

	project = applyAccessories(input, project)

	if err = models.UpdateProject(h.DB, project); err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update project")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, project)
}

// DeleteProject deletes a project
func (h *AdminHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Check if project exists
	project, err := models.GetProjectByID(h.DB, id)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve project")
		return
	}

	if project == nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	// Delete project
	if err := models.DeleteProject(h.DB, id); err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Project deleted successfully"})
}

func applyAccessories(input ProjectInput, project *models.Project) *models.Project {
	var tools []models.ProjectAccessory
	var categories []models.ProjectAccessory
	var skills []models.ProjectAccessory
	var supplies []models.ProjectAccessory
	var ages []models.ProjectAccessory

	for _, val := range input.Tools {
		a := models.ProjectAccessory{
			ID: val,
		}
		tools = append(tools, a)
	}
	project.Tools = tools

	for _, val := range input.Categories {
		a := models.ProjectAccessory{
			ID: val,
		}
		categories = append(categories, a)
	}
	project.Categories = categories

	for _, val := range input.Skills {
		a := models.ProjectAccessory{
			ID: val,
		}
		skills = append(skills, a)
	}
	project.Skills = skills

	for _, val := range input.Supplies {
		a := models.ProjectAccessory{
			ID: val,
		}
		supplies = append(supplies, a)
	}
	project.Supplies = supplies

	for _, val := range input.Ages {
		a := models.ProjectAccessory{
			ID: val,
		}
		ages = append(ages, a)
	}
	project.Ages = ages

	return project
}
