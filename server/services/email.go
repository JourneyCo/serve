package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-gomail/gomail"
	"serve/config"
	"serve/models"
)

const (
	OneDay       = "templates/one_day.html"
	OneWeek      = "templates/one_week.html"
	Registration = "registration.html"
	TwoWeeks     = "templates/two_week.html"
)

// EmailService handles email operations
type EmailService struct {
	Config *config.Config
	auth   smtp.Auth
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
func (s *EmailService) SendRegistrationConfirmation(user *models.User, project *models.Project) error {
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
	}{
		Name:            fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		ProjectTitle:    project.Title,
		ProjectDesc:     project.Description,
		Area:            project.Area,
		Address:         project.LocationAddress,
		ProjectDate:     projectDateFormatted,
		Time:            project.Time,
		ProjectDateFull: project.ProjectDate,
	}

	// Send the email
	return s.sendEmail(user.Email, subject, Registration, data)
}

// SendReminderEmail sends a reminder email for an upcoming project
func (s *EmailService) SendReminderEmail(registration *models.Registration, daysLeft int) error {
	var subject string
	var templateStr string

	// Set subject and template based on days left
	switch daysLeft {
	case 14:
		subject = fmt.Sprintf("2 Weeks Until Your Project: %s", registration.Project.Title)
		templateStr = TwoWeeks
	case 7:
		subject = fmt.Sprintf("1 Week Until Your Project: %s", registration.Project.Title)
		templateStr = OneWeek
	case 1:
		subject = fmt.Sprintf("Tomorrow: Your Project %s Begins", registration.Project.Title)
		templateStr = OneDay
	default:
		return fmt.Errorf("unsupported reminder interval: %d days", daysLeft)
	}

	// Format dates
	projectDateFormatted := registration.Project.ProjectDate.Format("Monday, January 2, 2006")

	// Create email data
	data := struct {
		Name         string
		ProjectTitle string
		ProjectDesc  string
		ProjectDate  string
		Time         string
		DaysLeft     int
	}{
		Name:         fmt.Sprintf("%s %s", registration.User.FirstName, registration.User.LastName),
		ProjectTitle: registration.Project.Title,
		ProjectDesc:  registration.Project.Description,
		ProjectDate:  projectDateFormatted,
		Time:         registration.Project.Time,
		DaysLeft:     daysLeft,
	}

	// Send the email
	return s.sendEmail(registration.User.Email, subject, templateStr, data)
}

// sendEmail is a helper function to send emails
func (s *EmailService) sendEmail(to, subject, templateStr string, data interface{}) error {
	// Parse template
	p := filepath.Join("templates", templateStr)
	name := path.Base(p)
	t, err := template.New(name).ParseFiles(p)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Execute template
	var body bytes.Buffer
	if err = t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", s.Config.MailFrom)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	// Set the plain-text version of the email
	message.SetBody(
		"text/plain",
		"This is a Test Email\n\nHello!\nThis is a test email with plain-text formatting.\nThanks,\nMailtrap",
	)

	// Set the HTML version of the email
	message.AddAlternative("text/html", body.String())
	port, err := strconv.Atoi(s.Config.MailPort)
	if err != nil {
		log.Println("error parsing mail port: ", err)
		return err
	}

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(s.Config.MailHost, port, s.Config.MailUser, s.Config.MailPass)

	// Send the email
	if err = dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		return err
	} else {
		fmt.Println("HTML Email sent successfully with a plain-text alternative!")
	}
	log.Printf("Email sent to %s: %s", to, subject)
	return nil
}
