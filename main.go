package main

import (
	"encoding/json"
	"fmt"
	"github.com/BohdanMelnyk/bus-shulter-checker/notification"
	"github.com/BohdanMelnyk/bus-shulter-checker/scheduler"
	"github.com/BohdanMelnyk/bus-shulter-checker/shuttle"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type CheckResult struct {
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	Available      bool      `json:"available"`
	CheckedDates   []string `json:"checkedDates"`
	AvailableDates []string `json:"availableDates"`
	CheckedAt      time.Time `json:"checkedAt"`
}

type AllChecksResponse struct {
	Results []CheckResult `json:"results"`
}

func main() {
	log.Println("Starting bus-shuttle-checker...")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Get configuration from environment variables
	mailgunDomain := os.Getenv("MAILGUN_DOMAIN")
	mailgunAPIKey := os.Getenv("MAILGUN_API_KEY")
	recipientEmail := os.Getenv("RECIPIENT_EMAIL")
	senderEmail := os.Getenv("SENDER_EMAIL")

	// Validate required environment variables
	var missingVars []string

	if mailgunDomain == "" {
		missingVars = append(missingVars, "MAILGUN_DOMAIN")
	}
	if mailgunAPIKey == "" {
		missingVars = append(missingVars, "MAILGUN_API_KEY")
	}
	if recipientEmail == "" {
		missingVars = append(missingVars, "RECIPIENT_EMAIL")
	}
	if senderEmail == "" {
		missingVars = append(missingVars, "SENDER_EMAIL")
	}

	if len(missingVars) > 0 {
		log.Fatalf("Missing required environment variables: %s", strings.Join(missingVars, ", "))
	}

	// Create an availability checker
	availabilityChecker := shuttle.NewWebChecker()

	// Create an email notifier
	emailNotifier := notification.NewEmailNotifier(
		mailgunDomain,
		mailgunAPIKey,
		recipientEmail,
		senderEmail,
	)

	// Create a scheduler with the checker and notifier
	availabilityScheduler := scheduler.NewAvailabilityScheduler(
		availabilityChecker,
		emailNotifier,
		1*time.Hour,
	)

	// Set up routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a closure to pass the emailNotifier to checkAllHandler
	http.HandleFunc("/check-all", func(w http.ResponseWriter, r *http.Request) {
		checkAllHandler(w, r, emailNotifier)
	})

	// Start periodic checks in a goroutine
	go availabilityScheduler.StartPeriodicChecks()

	// Start HTTP server for health checks
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting HTTP server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func checkAllHandler(w http.ResponseWriter, r *http.Request, notifier notification.Notifier) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var results []CheckResult
	checker := shuttle.NewWebChecker()

	log.Println("Starting availability check for all locations...")

	// Create a combined map of all shuttle locations to check
	combinedLocations := make(map[string]shuttle.ShuttleInfo)

	// Add locations from ShuttleURLs
	for name, location := range shuttle.ShuttleURLs {
		combinedLocations[name] = location
	}

	// Add additional locations
	additionalLocations := map[string]shuttle.ShuttleInfo{
		"Lake Louise Late Sep": {
			URL:   "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-09-30&endDate=2025-10-01&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T17:04:11.889&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
			Dates: []string{"2025-09-30", "2025-10-01"},
		},
		"Lake Louise Early Oct": {
			URL:   "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-10-07&endDate=2025-10-08&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,1,null%5D%5D&searchTime=2025-06-01T14:02:43.225&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
			Dates: []string{"2025-10-12"},
		},
	}

	// Add additional locations to combined map
	for name, location := range additionalLocations {
		combinedLocations[name] = location
	}

	// Check all locations
	log.Printf("Checking %d total locations...", len(combinedLocations))
	for name, location := range combinedLocations {
		log.Printf("Checking %s...", name)
		available, availableDates := checker.CheckAvailabilityForDates(location.URL, location.Dates)
		result := CheckResult{
			Name:           name,
			URL:            location.URL,
			Available:      available,
			CheckedDates:   location.Dates,
			AvailableDates: availableDates,
			CheckedAt:      time.Now(),
		}
		results = append(results, result)
		log.Printf("Check result for %s: %v (Available dates: %v)", name, available, availableDates)

		if available {
			message := fmt.Sprintf("Slots available for %s on dates: %s", name, strings.Join(availableDates, ", "))
			log.Printf("%s, sending notification...", message)
			if id, err := notifier.SendNotification(location.URL, message); err != nil {
				log.Printf("Error sending notification for %s: %v", name, err)
			} else {
				log.Printf("Notification sent successfully for %s, ID: %s", name, id)
			}
		}
	}

	log.Println("Availability check completed")

	response := AllChecksResponse{
		Results: results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
