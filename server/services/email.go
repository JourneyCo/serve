package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"time"

	"project-registration-system/config"
	"project-registration-system/models"
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
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPHost,
	)

	return &EmailService{
		Config: cfg,
		auth:   auth,
	}
}

// SendRegistrationConfirmation sends a confirmation email when a user registers for a project
func (s *EmailService) SendRegistrationConfirmation(user *models.User, project *models.Project) error {
	subject := fmt.Sprintf("Registration Confirmation: %s", project.Title)

	// Format dates
	projectDateFormatted := project.ProjectDate.Format("Monday, January 2, 2006")

	// Create email data
	data := struct {
		Name            string
		ProjectTitle    string
		ProjectDesc     string
		ProjectDate     string
		StartTime       string
		EndTime         string
		ProjectDateFull time.Time
	}{
		Name:            fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		ProjectTitle:    project.Title,
		ProjectDesc:     project.Description,
		ProjectDate:     projectDateFormatted,
		StartTime:       project.StartTime,
		EndTime:         project.EndTime,
		ProjectDateFull: project.ProjectDate,
	}

	// Send the email
	return s.sendEmail(user.Email, subject, registrationTemplate, data)
}

// SendReminderEmail sends a reminder email for an upcoming project
func (s *EmailService) SendReminderEmail(registration *models.Registration, daysLeft int) error {
	var subject string
	var templateStr string

	// Set subject and template based on days left
	switch daysLeft {
	case 14:
		subject = fmt.Sprintf("2 Weeks Until Your Project: %s", registration.Project.Title)
		templateStr = twoWeekReminderTemplate
	case 7:
		subject = fmt.Sprintf("1 Week Until Your Project: %s", registration.Project.Title)
		templateStr = oneWeekReminderTemplate
	case 1:
		subject = fmt.Sprintf("Tomorrow: Your Project %s Begins", registration.Project.Title)
		templateStr = oneDayReminderTemplate
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
		StartTime    string
		EndTime      string
		DaysLeft     int
	}{
		Name:         fmt.Sprintf("%s %s", registration.User.FirstName, registration.User.LastName),
		ProjectTitle: registration.Project.Title,
		ProjectDesc:  registration.Project.Description,
		ProjectDate:  projectDateFormatted,
		StartTime:    registration.Project.StartTime,
		EndTime:      registration.Project.EndTime,
		DaysLeft:     daysLeft,
	}

	// Send the email
	return s.sendEmail(registration.User.Email, subject, templateStr, data)
}

// sendEmail is a helper function to send emails
func (s *EmailService) sendEmail(to, subject, templateStr string, data interface{}) error {
	// Parse template
	t, err := template.New("email").Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Execute template
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Set up email headers
	headers := make(map[string]string)
	headers["From"] = s.Config.EmailFrom
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Create message
	var message bytes.Buffer
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.Write(body.Bytes())

	// Send the email
	addr := fmt.Sprintf("%s:%s", s.Config.SMTPHost, s.Config.SMTPPort)
	if err := smtp.SendMail(addr, s.auth, s.Config.EmailFrom, []string{to}, message.Bytes()); err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent to %s: %s", to, subject)
	return nil
}

// Email templates
const registrationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Registration Confirmation</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #3f51b5; color: white; padding: 15px; text-align: center; }
        .content { padding: 20px; border: 1px solid #ddd; }
        .footer { text-align: center; margin-top: 20px; color: #888; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Registration Confirmation</h1>
        </div>
        <div class="content">
            <p>Hello {{.Name}},</p>
            <p>Thank you for registering for the project: <strong>{{.ProjectTitle}}</strong>.</p>
            <p>Project Details:</p>
            <ul>
                <li><strong>Project:</strong> {{.ProjectTitle}}</li>
                <li><strong>Description:</strong> {{.ProjectDesc}}</li>
                <li><strong>Date:</strong> {{.ProjectDate}}</li>
                <li><strong>Time:</strong> {{.StartTime}} - {{.EndTime}}</li>
            </ul>
            <p>We'll send you reminder emails as the project date approaches.</p>
            <p>Please contact us if you have any questions or need to make changes to your registration.</p>
            <p>Thank you,<br>The Project Registration Team</p>
        </div>
        <div class="footer">
            <p>This is an automated message, please do not reply.</p>
        </div>
    </div>
</body>
</html>
`

const twoWeekReminderTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Project Reminder - 2 Weeks</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #3f51b5; color: white; padding: 15px; text-align: center; }
        .content { padding: 20px; border: 1px solid #ddd; }
        .footer { text-align: center; margin-top: 20px; color: #888; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Project Reminder - 2 Weeks</h1>
        </div>
        <div class="content">
            <p>Hello {{.Name}},</p>
            <p>This is a friendly reminder that your registered project <strong>{{.ProjectTitle}}</strong> starts in 2 weeks on {{.ProjectDate}} at {{.StartTime}}.</p>
            <p>Project Details:</p>
            <ul>
                <li><strong>Project:</strong> {{.ProjectTitle}}</li>
                <li><strong>Description:</strong> {{.ProjectDesc}}</li>
                <li><strong>Date:</strong> {{.ProjectDate}}</li>
                <li><strong>Time:</strong> {{.StartTime}} - {{.EndTime}}</li>
            </ul>
            <p>Please ensure you have made all necessary preparations for the project.</p>
            <p>If you need to cancel your registration, please do so at least 3 days before the project starts.</p>
            <p>Thank you,<br>The Project Registration Team</p>
        </div>
        <div class="footer">
            <p>This is an automated message, please do not reply.</p>
        </div>
    </div>
</body>
</html>
`

const oneWeekReminderTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Project Reminder - 1 Week</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #ff9800; color: white; padding: 15px; text-align: center; }
        .content { padding: 20px; border: 1px solid #ddd; }
        .footer { text-align: center; margin-top: 20px; color: #888; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Project Reminder - 1 Week</h1>
        </div>
        <div class="content">
            <p>Hello {{.Name}},</p>
            <p>This is a reminder that your registered project <strong>{{.ProjectTitle}}</strong> starts in 1 week on {{.ProjectDate}} at {{.StartTime}}.</p>
            <p>Project Details:</p>
            <ul>
                <li><strong>Project:</strong> {{.ProjectTitle}}</li>
                <li><strong>Description:</strong> {{.ProjectDesc}}</li>
                <li><strong>Date:</strong> {{.ProjectDate}}</li>
                <li><strong>Time:</strong> {{.StartTime}} - {{.EndTime}}</li>
            </ul>
            <p>Please make sure you are prepared for the project date.</p>
            <p>If you need to cancel your registration, please do so as soon as possible.</p>
            <p>Thank you,<br>The Project Registration Team</p>
        </div>
        <div class="footer">
            <p>This is an automated message, please do not reply.</p>
        </div>
    </div>
</body>
</html>
`

const oneDayReminderTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Project Reminder - Tomorrow</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f44336; color: white; padding: 15px; text-align: center; }
        .content { padding: 20px; border: 1px solid #ddd; }
        .footer { text-align: center; margin-top: 20px; color: #888; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Project Starts Tomorrow!</h1>
        </div>
        <div class="content">
            <p>Hello {{.Name}},</p>
            <p>This is your final reminder that your registered project <strong>{{.ProjectTitle}}</strong> starts tomorrow ({{.ProjectDate}}) at {{.StartTime}}.</p>
            <p>Project Details:</p>
            <ul>
                <li><strong>Project:</strong> {{.ProjectTitle}}</li>
                <li><strong>Description:</strong> {{.ProjectDesc}}</li>
                <li><strong>Date:</strong> {{.ProjectDate}}</li>
                <li><strong>Time:</strong> {{.StartTime}} - {{.EndTime}}</li>
            </ul>
            <p>Please ensure you are fully prepared and ready to begin.</p>
            <p>We look forward to your participation!</p>
            <p>Thank you,<br>The Project Registration Team</p>
        </div>
        <div class="footer">
            <p>This is an automated message, please do not reply.</p>
        </div>
    </div>
</body>
</html>
`
