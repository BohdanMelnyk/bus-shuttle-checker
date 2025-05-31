package scheduler

import (
	"fmt"
	"github.com/BohdanMelnyk/bus-shulter-checker/notification"
	"github.com/BohdanMelnyk/bus-shulter-checker/shuttle"
	"time"
)

// AvailabilityScheduler manages periodic checking for shuttle availability
type AvailabilityScheduler struct {
	checker  shuttle.AvailabilityChecker
	notifier notification.Notifier
	interval time.Duration
}

// NewAvailabilityScheduler creates a new scheduler with the specified checker, notifier, and interval
func NewAvailabilityScheduler(checker shuttle.AvailabilityChecker, notifier notification.Notifier, interval time.Duration) *AvailabilityScheduler {
	return &AvailabilityScheduler{
		checker:  checker,
		notifier: notifier,
		interval: interval,
	}
}

// CheckAvailabilityAndNotify runs a single check across all URLs and sends notifications
// for any available slots
func (s *AvailabilityScheduler) CheckAvailabilityAndNotify() {
	fmt.Println("Starting availability check...")
	for name, url := range shuttle.GetAllURLs() {
		if s.checker.CheckAvailability(url) {
			fmt.Printf("Slot available for %s: %s\n", name, url)
			id, err := s.notifier.SendNotification(url, name)
			if err != nil {
				fmt.Printf("Error sending notification for %s %s: %v\n", name, url, err)
			} else {
				fmt.Printf("Notification sent successfully for %s %s, ID: %s\n", name, url, id)
			}
		} else {
			fmt.Printf("No slots available for %s: %s\n", name, url)
		}
	}
	fmt.Println("Availability check completed")
}

// StartPeriodicChecks begins a periodic check for availability at the specified interval
// This function blocks and runs indefinitely
func (s *AvailabilityScheduler) StartPeriodicChecks() {
	fmt.Println("Starting periodic checks for shuttle availability...")

	// Run immediately on startup
	s.CheckAvailabilityAndNotify()

	// Then run on the specified interval
	for range time.Tick(s.interval) {
		s.CheckAvailabilityAndNotify()
	}
}
