package services

import (
	"database/sql"
	"log"
	"time"

	"serve/models"
)

// Scheduler handles scheduling of email reminders
type Scheduler struct {
	DB           *sql.DB
	EmailService *EmailService
	TextService  *TextService
	stop         chan struct{}
}

// NewScheduler creates a new scheduler service
func NewScheduler(db *sql.DB, emailService *EmailService, textService *TextService) *Scheduler {
	return &Scheduler{
		DB:           db,
		EmailService: emailService,
		TextService:  textService,
		stop:         make(chan struct{}),
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	log.Println("Starting email reminder scheduler...")

	// Run immediately on startup
	s.processReminders()

	// Set up a ticker to run daily at a specific time (e.g., 8:00 AM)
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.processReminders()
		case <-s.stop:
			log.Println("Stopping email reminder scheduler...")
			return
		}
	}
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	close(s.stop)
}

// processReminders processes all reminders
func (s *Scheduler) processReminders() {
	log.Println("Processing email reminders...")

	// Process reminders for different intervals
	s.processReminderForDays(14) // 2 weeks before
	s.processReminderForDays(7)  // 1 week before
	// s.processReminderForDays(1)  // 1 day before // NOT using per Cory

	log.Println("Finished processing email reminders")
}

// processReminderForDays processes reminders for a specific day interval
func (s *Scheduler) processReminderForDays(days int) {
	registrations, err := models.GetRegistrationsForReminders(s.DB, days)
	if err != nil {
		log.Printf("Error getting registrations for %d days reminder: %v", days, err)
		return
	}

	log.Printf("Found %d registrations for %d days reminder", len(registrations), days)

	// send emails - we are rate limited to 200 emails per hour by mailtrap, so we will limit ourselves to 150
	// just to be safe. This is in case additional people register or obtain emails in the hour while we are
	// sending registration emails
	ticker := time.NewTicker(24 * time.Second) // 3600/180 = 24s per email
	defer ticker.Stop()

	for _, reg := range registrations {
		if err := s.EmailService.SendReminderEmail(&reg, days); err != nil {
			log.Printf("Error sending %d days reminder email to %s: %v", days, reg.User.Email, err)
		}
		<-ticker.C // Wait for the next tick before sending the next email
	}

	// send text messages - not doing this in 2025
	// var list []models.Registration
	// for _, reg := range registrations {
	// 	if reg.User.TextPermission { // exclude users who do not want texts
	// 		list = append(list, reg)
	// 	}
	// }
	// if err := s.TextService.SendReminderText(list, days); err != nil {
	// 	log.Printf(
	// 		"Error sending %d days reminder text: %v", days, err,
	// 	)
	// }
}
