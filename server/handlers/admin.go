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
	GoogleID             *int    `json:"google_id"`
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
	LocationAddress      string  `json:"location_address"`
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
}

// RegisterAdminRoutes registers the routes for admin handlers
func RegisterAdminRoutes(router *mux.Router, db *sql.DB) {
	handler := &AdminHandler{
		DB: db,
	}

	router.HandleFunc("/users", handler.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/registrations", handler.GetAllRegistrations).Methods(http.MethodGet)
	router.HandleFunc("/projects", handler.CreateProject).Methods(http.MethodPost)
	router.HandleFunc("/projects/{id:[0-9]+}", handler.UpdateProject).Methods(http.MethodPut)
	router.HandleFunc("/projects/{id:[0-9]+}", handler.DeleteProject).Methods(http.MethodDelete)
	router.HandleFunc("/registrations/{id:[0-9]+}", handler.DeleteRegistration).Methods(http.MethodDelete)
}

// GetAllRegistrations returns all registrations across all projects
func (h *AdminHandler) GetAllRegistrations(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT r.id, r.user_id, r.project_id, r.status, r.guest_count, r.lead_interest,
		r.created_at, r.updated_at,
		u.email, u.first_name, u.last_name,
		p.title, p.description, p.time, p.project_date
		FROM registrations r
		JOIN users u ON r.user_id = u.id
		JOIN projects p ON r.project_id = p.id
		ORDER BY r.created_at DESC
	`

	rows, err := h.DB.Query(query)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve registrations")
		return
	}
	defer rows.Close()

	var registrations []models.Registration
	for rows.Next() {
		var r models.Registration
		r.User = &models.User{}
		r.Project = &models.Project{}

		err := rows.Scan(
			&r.ID, &r.UserID, &r.ProjectID, &r.Status, &r.GuestCount, &r.LeadInterest,
			&r.CreatedAt, &r.UpdatedAt,
			&r.User.Email, &r.User.FirstName, &r.User.LastName,
			&r.Project.Title, &r.Project.Description, &r.Project.Time, &r.Project.ProjectDate,
		)
		if err != nil {
			middleware.RespondWithError(w, http.StatusInternalServerError, "Error scanning registrations")
			return
		}

		r.User.ID = r.UserID
		r.Project.ID = r.ProjectID
		registrations = append(registrations, r)
	}

	middleware.RespondWithJSON(w, http.StatusOK, registrations)
}

// UpdateRegistrationGuestCount updates the guest count for a registration
func (h *AdminHandler) UpdateRegistrationGuestCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regID, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid registration ID")
		return
	}

	var input struct {
		GuestCount int `json:"guest_count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.GuestCount < 0 {
		middleware.RespondWithError(w, http.StatusBadRequest, "Guest count cannot be negative")
		return
	}

	query := `UPDATE registrations SET guest_count = $1 WHERE id = $2`
	result, err := h.DB.Exec(query, input.GuestCount, regID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update registration")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		middleware.RespondWithError(w, http.StatusNotFound, "Registration not found")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Registration updated successfully"})
}

// DeleteRegistration deletes a registration
func (h *AdminHandler) DeleteRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regID, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid registration ID")
		return
	}

	query := `DELETE FROM registrations WHERE id = $1`
	result, err := h.DB.Exec(query, regID)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to delete registration")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		middleware.RespondWithError(w, http.StatusNotFound, "Registration not found")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Registration deleted successfully"})
}

// GetAllUsers returns all users
func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := models.GetAllUsers(ctx, h.DB)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, users)
}

// CreateProject creates a new project
func (h *AdminHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	projectDate, err := time.Parse(time.RFC3339, input.ProjectDate)

	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project date format")
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
		LocationAddress:  input.LocationAddress,
		Latitude:         input.Latitude,
		Longitude:        input.Longitude,
		LeadUserID:       input.LeadUserID,
	}

	project = applyAccessories(input, project)

	if err = models.CreateProject(ctx, h.DB, project); err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	middleware.RespondWithJSON(w, http.StatusCreated, project)
}

// UpdateProject updates an existing project
func (h *AdminHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Check if project exists
	project, err := models.GetProjectByID(ctx, h.DB, id)
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

	projectDate, err := time.Parse(time.RFC3339, input.ProjectDate)
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project date format (use YYYY-MM-DD)")
		return
	}

	// Update project
	project.GoogleID = input.GoogleID
	project.Title = input.Title
	project.Description = input.Description
	project.ShortDescription = input.ShortDescription
	project.Time = input.Time
	project.ProjectDate = projectDate
	project.MaxCapacity = input.MaxCapacity
	project.LocationName = input.LocationName
	project.LocationAddress = input.LocationAddress
	project.Latitude = input.Latitude
	project.Longitude = input.Longitude
	project.WheelchairAccessible = input.WheelchairAccessible

	// TODO: Need to delete accessories here for an existing project before adding accessories

	project = applyAccessories(input, project)

	if err = models.UpdateProject(ctx, h.DB, project); err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to update project")
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, project)
}

// DeleteProject deletes a project
func (h *AdminHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("invalid project id")
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Check if project exists
	project, err := models.GetProjectByID(ctx, h.DB, id)
	if err != nil {
		log.Println("error querying for project")
		middleware.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve project")
		return
	}

	if project == nil {
		middleware.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	// Delete project
	if err := models.DeleteProject(ctx, h.DB, id); err != nil {
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
