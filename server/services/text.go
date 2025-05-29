package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"time"

	"serve/config"
	"serve/models"
)

const (
	clearstreamTextURL  = "https://api.getclearstream.com/v1/texts"
	clearstreamTextFrom = "94000"
)

// TextService handles text operations
type TextService struct {
	APIKey string
	Config *config.Config
}

type ClearStreamRequest struct {
	To         []string              `json:"to"`
	From       string                `json:"from"`
	TextHeader string                `json:"text_header"`
	TextBody   string                `json:"text_body"`
	List       []models.Registration `json:"-"`
	APIKey     string                `json:"-"`
	PhoneList  []string              `json:"-"`

	// DefaultHeader is a flag to use the system's default header for journey
	DefaultHeader bool `json:"use_default_header"`

	// OverRideOptOut is used to override a subscribers wish to opt out of
	// texts. This should never be set to true per Journey's wishes.
	OverRideOptOut bool `json:"override_optouts"`
}

type ClearStreamResponse struct {
	Data struct {
		ID       any       `json:"id"`
		Status   string    `json:"status"`
		QueuedAt time.Time `json:"queued_at"`
		Text     string    `json:"text"`
		To       []string  `json:"to"`
		Skipped  []any     `json:"skipped"`
		From     string    `json:"from"`
		Media    []any     `json:"media"`
	} `json:"data"`
}

// NewTextService creates a new text service
func NewTextService(cfg *config.Config) *TextService {
	return &TextService{
		APIKey: cfg.ClearStreamAPIKey,
		Config: cfg,
	}
}

// SendRegistrationConfirmation sends a confirmation text when a user registers for a project
func (s *TextService) SendRegistrationConfirmation(user *models.User, project *models.Project) error {
	subject := fmt.Sprintf("Registration Confirmation: %s", project.Title)

	// Format dates
	// projectDateFormatted := project.ProjectDate.Format("Monday, January 2, 2006")

	r := models.Registration{
		User: user,
	}

	// Create text data
	req := ClearStreamRequest{
		From:       clearstreamTextFrom,
		TextHeader: "Journey Serve Day",
		TextBody:   subject,
		List:       []models.Registration{r},
		APIKey:     s.APIKey,
	}

	// Send the text
	return req.sendText()
}

// SendReminderText sends a reminder text for an upcoming project
func (s *TextService) SendReminderText(list []models.Registration, daysLeft int) error {
	var subject string

	if len(list) == 0 {
		log.Println(fmt.Sprintf("Found 0 registrations for %d days reminder texts", daysLeft))
		return nil
	}

	// Set subject and template based on days left
	switch daysLeft {
	case 14:
		subject = fmt.Sprintf("2 Weeks Until Your Serve Project!")
	case 7:
		subject = fmt.Sprintf("1 Week Until Your Serve Project!")
	case 1:
		subject = fmt.Sprintf("Tomorrow: Your Serve Project Begins!")
	default:
		return fmt.Errorf("unsupported reminder interval: %d days", daysLeft)
	}

	// filter out registrations that do not want to be sent texts
	allowedList := slices.DeleteFunc(
		list, func(r models.Registration) bool {
			return r.User.TextPermission
		},
	)

	req := ClearStreamRequest{
		From:       clearstreamTextFrom,
		TextHeader: "Journey Serve Day",
		TextBody:   subject,
		List:       allowedList,
		APIKey:     s.APIKey,
	}

	return req.sendText()
}

func (s *TextService) SendTestText() error {
	req := ClearStreamRequest{
		From:       clearstreamTextFrom,
		TextHeader: "Journey Serve Day",
		TextBody:   "Journey Serve",
		List: []models.Registration{
			{
				ID: 1,
				User: &models.User{
					Phone: "+13039477791",
				},
			},
		},
		APIKey: s.APIKey,
	}

	return req.sendText()

}

// sendText is a helper function to send texts
func (c *ClearStreamRequest) sendText() error {
	var sendList []string
	for _, num := range c.List {
		sendList = append(sendList, num.User.Phone)
	}

	if len(sendList) == 0 {
		log.Println("no phone numbers to send text reminder(s) to")
		return nil
	}

	csr := ClearStreamRequest{
		To:            sendList,
		From:          c.From,
		TextHeader:    "Journey Serve Day",
		TextBody:      c.TextBody,
		DefaultHeader: true,
		APIKey:        c.APIKey,
	}

	b, err := json.Marshal(csr)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(http.MethodPost, clearstreamTextURL, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("X-API-KEY", c.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post text request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	log.Println("text sent")

	var csResp ClearStreamResponse
	if err = json.Unmarshal(body, &csResp); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	log.Printf("%v", csResp)
	log.Printf("texts skipped: %v", csResp.Data.Skipped)
	return nil
}
