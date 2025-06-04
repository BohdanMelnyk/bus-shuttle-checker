package main

import (
	"encoding/json"
	"fmt"
	"github.com/BohdanMelnyk/bus-shulter-checker/notification"
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

type dummyResponseWriter struct{}

func (d *dummyResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (d *dummyResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (d *dummyResponseWriter) WriteHeader(statusCode int) {}

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

	// Create an API client
	apiClient := shuttle.NewAPIClient()

	// Create an email notifier
	emailNotifier := notification.NewEmailNotifier(
		mailgunDomain,
		mailgunAPIKey,
		recipientEmail,
		senderEmail,
	)

	// Set up routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a closure to pass the emailNotifier to checkAllHandler
	http.HandleFunc("/check-all", func(w http.ResponseWriter, r *http.Request) {
		checkAllHandler(w, r, emailNotifier, apiClient)
	})

	// Create a channel to signal shutdown
	shutdownChan := make(chan struct{})

	// Start HTTP server in a goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Starting HTTP server on port %s...\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Printf("HTTP server stopped: %v\n", err)
		}
	}()

	// Run one check immediately
	log.Println("Running initial availability check...")
	checkAllLocations(emailNotifier, apiClient)

	// Set a timer for 5 minutes
	shutdownTimer := time.NewTimer(5 * time.Minute)

	// Wait for either shutdown signal or timer expiration
	select {
	case <-shutdownTimer.C:
		log.Println("5 minutes elapsed, shutting down...")
	case <-shutdownChan:
		log.Println("Received shutdown signal...")
	}

	log.Println("Shutdown complete")
	os.Exit(0)
}

func checkAllHandler(w http.ResponseWriter, r *http.Request, notifier notification.Notifier, apiClient *shuttle.APIClient) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var results []CheckResult
	log.Println("Starting availability check for all locations...")

	for _, location := range shuttle.Locations {
		log.Printf("Checking %s...", location.Name)
		
		// Get the first and last date to check
		if len(location.Dates) == 0 {
			continue
		}
		startDate := location.Dates[0]
		endDate := location.Dates[len(location.Dates)-1]

		available, availableDates, err := apiClient.HasAvailability(
			location.Name,
			location.LocationID,
			startDate,
			endDate,
			location.ResourceIDs,
			location.BookingCategory,
		)
		if err != nil {
			log.Printf("Error checking availability for %s: %v", location.Name, err)
			continue
		}

		result := CheckResult{
			Name:           location.Name,
			URL:            fmt.Sprintf("https://reservation.pc.gc.ca/create-booking/results?resourceLocationId=%d", location.LocationID),
			Available:      available,
			CheckedDates:   location.Dates,
			AvailableDates: availableDates,
			CheckedAt:      time.Now(),
		}
		results = append(results, result)
		log.Printf("Check result for %s: %v (Available dates: %v)", location.Name, available, availableDates)

		if available {
			message := fmt.Sprintf("Slots available for %s on dates: %s", location.Name, strings.Join(availableDates, ", "))
			log.Printf("%s, sending notification...", message)
			if id, err := notifier.SendNotification(result.URL, message); err != nil {
				log.Printf("Error sending notification for %s: %v", location.Name, err)
			} else {
				log.Printf("Notification sent successfully for %s, ID: %s", location.Name, id)
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

func checkAllLocations(notifier notification.Notifier, apiClient *shuttle.APIClient) {
	checkAllHandler(&dummyResponseWriter{}, &http.Request{Method: http.MethodGet}, notifier, apiClient)
}
