package services

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"path"
	"path/filepath"
	"time"

	"serve/config"
	"serve/models"
)

const (
	OneDay       = "one_day.html"
	OneWeek      = "one_week.html"
	Registration = "registration.html"
	ThankYou     = "thank_you.html"
	TwoWeeks     = "two_week.html"
)

// EmailService handles email operations
type EmailService struct {
	Config *config.Config
	auth   smtp.Auth
}

// EmailService handles email operations
type mailtrapResponse struct {
	Success bool `json:"success"`
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) *EmailService {
	auth := smtp.PlainAuth(
		"",
		cfg.MailUser,
		cfg.MailPass,
		cfg.MailHost,
	)

	return &EmailService{
		Config: cfg,
		auth:   auth,
	}
}

// SendRegistrationConfirmation sends a confirmation email when a user registers for a project
func (s *EmailService) SendRegistrationConfirmation(user *models.User, project *models.Project, guests int) {
	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Hour)
	defer cancel()
	subject := fmt.Sprintf("Serve Day Project Confirmation")

	// Format dates
	projectDateFormatted := project.ProjectDate.Format("Monday, January 2, 2006")

	// Create email data
	data := struct {
		Name            string
		ProjectTitle    string
		ProjectDesc     string
		Area            string
		Address         string
		ProjectDate     string
		Time            string
		ProjectDateFull time.Time
		Guests          int
	}{
		Name:            fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		ProjectTitle:    project.Title,
		ProjectDesc:     project.Description,
		Area:            project.Area,
		Address:         project.LocationAddress,
		ProjectDate:     projectDateFormatted,
		Time:            project.Time,
		ProjectDateFull: project.ProjectDate,
		Guests:          guests,
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for attempt := 1; ; attempt++ {
		select {
		case <-ctx.Done():
			log.Printf("Email failed after 24 hours: %v", data)
			return
		case <-ticker.C:
			err := s.sendEmail(ctx, user.Email, subject, Registration, data)
			if err == nil {
				log.Printf("Email succeeded on attempt %d to %s", attempt, user.Email)
				return
			}
			log.Printf("Email Attempt %d failed: %v", attempt, err)
		}
	}
}

// SendReminderEmail sends a reminder email for an upcoming project
func (s *EmailService) SendReminderEmail(registration *models.Registration, daysLeft int) error {
	var subject string
	var templateStr string

	// Set subject and template based on days left
	switch daysLeft {
	case 14:
		subject = fmt.Sprintf("2 Weeks Until Your Journey Serve Day Project: %s", registration.Project.Title)
		templateStr = TwoWeeks
	case 7:
		subject = fmt.Sprintf("1 Week Until Your Journey Serve Day Project: %s", registration.Project.Title)
		templateStr = OneWeek
	// case 1:
	// 	subject = fmt.Sprintf("Tomorrow: Your Project %s Begins", registration.Project.Title)
	// 	templateStr = OneDay
	default:
		return fmt.Errorf("unsupported reminder interval: %d days", daysLeft)
	}

	// Format dates
	projectDateFormatted := registration.Project.ProjectDate.Format("Monday, January 2, 2006")

	// Create email data
	data := struct {
		Name             string
		ProjectTitle     string
		ProjectDesc      string
		ProjectDate      string
		Time             string
		DaysLeft         int
		ServeLeaderName  string
		ServeLeaderEmail string
		Guests           int
		Area             string
		Address          string
	}{
		Name:             fmt.Sprintf("%s %s", registration.User.FirstName, registration.User.LastName),
		ProjectTitle:     registration.Project.Title,
		ProjectDesc:      registration.Project.Description,
		ProjectDate:      projectDateFormatted,
		Time:             registration.Project.Time,
		DaysLeft:         daysLeft,
		ServeLeaderEmail: registration.Project.ServeLeadEmail,
		ServeLeaderName:  registration.Project.ServeLeadName,
		Guests:           registration.GuestCount,
		Area:             registration.Project.Area,
		Address:          registration.Project.LocationAddress,
	}

	ctx := context.Background()
	// Send the email
	return s.sendEmail(ctx, registration.User.Email, subject, templateStr, data)
}

// sendEmail is a helper function to send emails
func (s *EmailService) sendEmail(ctx context.Context, to, subject, templateStr string, data interface{}) error {
	// Parse template
	p := filepath.Join("templates", templateStr)
	name := path.Base(p)
	t, err := template.New(name).ParseFiles(p)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Execute template
	var htmlBody bytes.Buffer
	if err = t.Execute(&htmlBody, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	payload := map[string]interface{}{
		"from":    map[string]string{"email": s.Config.MailFrom},
		"to":      []map[string]string{{"email": to}},
		"subject": subject,
		"text":    "You are confirmed for Serve Day",
		"html":    htmlBody.String(),
		"reply_to": map[string]string{
			"email": "sarawiest@journeycolorado.com",
			"name":  "Sara Wiest",
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email template: %w", err)
	}

	// Replace with your Mailtrap API token
	token := s.Config.MailKey
	// Mailtrap API endpoint
	apiURL := s.Config.MailHost

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://"+apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create email request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request
	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email via provider: %w", err)
	}
	defer res.Body.Close()

	// Read and print response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var mtr mailtrapResponse
	if err = json.Unmarshal(body, &mtr); err != nil {
		return fmt.Errorf("error unmarshaling response")
	}

	if !mtr.Success {
		return fmt.Errorf("email call completed but not successful: %v", string(body))
	}

	log.Printf("Email sent to %s: %s", to, subject)
	return nil

}

// SendThankYouToAllUsers sends thank-you emails to all users in the database
func (s *EmailService) SendThankYouToAllUsers(ctx context.Context, db *sql.DB) error {
	// Get all users from the database
	users, err := models.GetAllUsers(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	log.Printf("Sending thank you emails to %d users", len(users))

	subject := "Serve Day - Thank you"

	// Rate limit to 150 emails per hour (24 seconds between emails)
	ticker := time.NewTicker(24 * time.Second)
	defer ticker.Stop()

	for i, user := range users {
		// Create email data
		data := struct {
			Name string
		}{
			Name: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		}

		// Send email with retry logic
		err := s.sendEmailWithRetry(ctx, user.Email, subject, ThankYou, data)
		if err != nil {
			log.Printf("Failed to send thank you email to %s: %v", user.Email, err)
			continue
		}

		log.Printf("Sent thank you email %d/%d to %s", i+1, len(users), user.Email)

		// Wait for rate limit except for the last email
		if i < len(users)-1 {
			<-ticker.C
		}
	}

	log.Printf("Finished sending thank you emails")
	return nil
}

// sendEmailWithRetry sends an email with retry logic
func (s *EmailService) sendEmailWithRetry(
	ctx context.Context, to, subject, templateStr string, data interface{},
) error {
	maxRetries := 3
	baseDelay := 1 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := s.sendEmail(ctx, to, subject, templateStr, data)
		if err == nil {
			return nil
		}

		if attempt == maxRetries {
			return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
		}

		// Exponential backoff
		delay := baseDelay * time.Duration(1<<uint(attempt-1))
		log.Printf("Attempt %d failed for %s, retrying in %v: %v", attempt, to, delay, err)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return nil
}
