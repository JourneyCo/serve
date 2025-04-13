package services

import (
	"database/sql"
	"log"
	"time"

	"project-registration-system/models"
)

// Scheduler handles scheduling of email reminders
type Scheduler struct {
	DB           *sql.DB
	EmailService *EmailService
	stop         chan struct{}
}

// NewScheduler creates a new scheduler service
func NewScheduler(db *sql.DB, emailService *EmailService) *Scheduler {
	return &Scheduler{
		DB:           db,
		EmailService: emailService,
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
	s.processReminderForDays(1)  // 1 day before

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

	for _, reg := range registrations {
		if err := s.EmailService.SendReminderEmail(&reg, days); err != nil {
			log.Printf("Error sending %d days reminder email to %s: %v", days, reg.User.Email, err)
		}
	}
}
