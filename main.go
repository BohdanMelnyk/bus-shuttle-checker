package main

import (
	"encoding/json"
	"github.com/BohdanMelnyk/bus-shulter-checker/notification"
	"github.com/BohdanMelnyk/bus-shulter-checker/scheduler"
	"github.com/BohdanMelnyk/bus-shulter-checker/shuttle"
	"log"
	"net/http"
	"os"
	"time"
)

type CheckResult struct {
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Available bool      `json:"available"`
	CheckedAt time.Time `json:"checkedAt"`
}

type AllChecksResponse struct {
	Results []CheckResult `json:"results"`
}

func main() {
	log.Println("Starting bus-shuttle-checker...")

	// Get configuration from environment variables
	mailgunDomain := os.Getenv("MAILGUN_DOMAIN")
	mailgunAPIKey := os.Getenv("MAILGUN_API_KEY")
	recipientEmail := os.Getenv("RECIPIENT_EMAIL")
	senderEmail := os.Getenv("SENDER_EMAIL")

	// Validate required environment variables
	if mailgunDomain == "" || mailgunAPIKey == "" || recipientEmail == "" || senderEmail == "" {
		log.Fatal("Missing required environment variables. Please set MAILGUN_DOMAIN, MAILGUN_API_KEY, RECIPIENT_EMAIL, and SENDER_EMAIL")
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

	// Get all URLs from the shuttle package
	allURLs := shuttle.GetAllURLs()
	log.Printf("Checking %d locations from urls.go...", len(allURLs))

	// Check URLs from urls.go
	for name, url := range allURLs {
		log.Printf("Checking %s...", name)
		available := checker.CheckAvailability(url)
		result := CheckResult{
			Name:      name,
			URL:       url,
			Available: available,
			CheckedAt: time.Now(),
		}
		results = append(results, result)
		log.Printf("Check result for %s: %v", name, available)

		if available {
			log.Printf("Slots available for %s, sending notification...", name)
			if id, err := notifier.SendNotification(url, name); err != nil {
				log.Printf("Error sending notification for %s: %v", name, err)
			} else {
				log.Printf("Notification sent successfully for %s, ID: %s", name, id)
			}
		}
	}

	// Additional URLs to check (not in urls.go)
	additionalURLs := map[string]string{
		"Lake Louise Late Sep": "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-09-30&endDate=2025-10-01&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T17:04:11.889&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
	}

	// Check additional URLs
	log.Printf("Checking %d additional locations...", len(additionalURLs))
	for name, url := range additionalURLs {
		log.Printf("Checking %s...", name)
		available := checker.CheckAvailability(url)
		result := CheckResult{
			Name:      name,
			URL:       url,
			Available: available,
			CheckedAt: time.Now(),
		}
		results = append(results, result)
		log.Printf("Check result for %s: %v", name, available)

		if available {
			log.Printf("Slots available for %s, sending notification...", name)
			if id, err := notifier.SendNotification(url, name); err != nil {
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
