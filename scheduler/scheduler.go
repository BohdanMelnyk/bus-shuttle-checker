package scheduler

import (
	"fmt"
	"github.com/BohdanMelnyk/bus-shulter-checker/notification"
	"github.com/BohdanMelnyk/bus-shulter-checker/shuttle"
	"strings"
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
	for name, location := range shuttle.ShuttleURLs {
		hasAvailability, availableDates := s.checker.CheckAvailabilityForDates(location.URL, location.Dates)
		if hasAvailability {
			message := fmt.Sprintf("Slot available for %s on dates: %s", name, strings.Join(availableDates, ", "))
			fmt.Printf("%s\n", message)
			id, err := s.notifier.SendNotification(location.URL, message)
			if err != nil {
				fmt.Printf("Error sending notification for %s %s: %v\n", name, location.URL, err)
			} else {
				fmt.Printf("Notification sent successfully for %s %s, ID: %s\n", name, location.URL, id)
			}
		} else {
			fmt.Printf("No slots available for %s: %s (checked dates: %s)\n", name, location.URL, strings.Join(location.Dates, ", "))
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
