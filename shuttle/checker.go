package shuttle

import (
	"github.com/go-rod/rod"
	"os"
	"sort"
	"strings"
)

// WebChecker implements the AvailabilityChecker interface using web scraping
type WebChecker struct{}

// NewWebChecker creates a new WebChecker instance
func NewWebChecker() *WebChecker {
	return &WebChecker{}
}

// initBrowser initializes a new browser instance with appropriate settings
func (w *WebChecker) initBrowser() *rod.Browser {
	browser := rod.New()
	
	// Check if running in CI environment
	if os.Getenv("CI") == "true" {
		browser = browser.NoSandbox()
	}
	
	return browser.MustConnect()
}

// CheckAvailability checks if a shuttle slot is available for the given URL
// It returns true if a slot is available, false otherwise
func (w *WebChecker) CheckAvailability(url string) bool {
	// Start a Rod browser
	browser := w.initBrowser()
	defer browser.MustClose()

	// Navigate to the URL
	page := browser.MustPage(url)

	// Wait for the page to load completely
	page.MustWaitLoad()

	// Additional wait to ensure all dynamic content is loaded
	page.MustWaitIdle()

	// Wait for all asynchronous processes to complete
	page.MustWaitDOMStable()

	// Wait for the specific element indicating the page is fully loaded
	page.MustElement(".chart-cell").MustWaitVisible()

	// Extract the page source
	pageSource := page.MustHTML()

	// Check for availability based on aria-label containing "Available"
	available := strings.Contains(pageSource, "Departures Available")
	return available
}

// CheckAvailabilityForDates checks if a shuttle slot is available for the given URL and specific dates
// It returns true if a slot is available on any of the specified dates, false otherwise
// Additionally, it returns a list of unique dates for which slots are available
func (w *WebChecker) CheckAvailabilityForDates(url string, dates []string) (bool, []string) {
	// Start a Rod browser
	browser := w.initBrowser()
	defer browser.MustClose()

	// Navigate to the URL
	page := browser.MustPage(url)

	// Wait for the page to load completely
	page.MustWaitLoad()

	// Additional wait to ensure all dynamic content is loaded
	page.MustWaitIdle()

	// Wait for all asynchronous processes to complete
	page.MustWaitDOMStable()

	// Extract all elements with the class "chart-cell"
	cells := page.MustElements(".chart-cell")

	// Initialize a map to store unique available dates
	availableDatesMap := make(map[string]struct{})

	// Iterate over the cells to find matching dates and availability
	for _, cell := range cells {
		ariaLabel := cell.MustAttribute("aria-label")
		dataE2eDate := cell.MustAttribute("data-e2e-date")

		if ariaLabel != nil && dataE2eDate != nil {
			if strings.Contains(*ariaLabel, "Departures Available") {
				for _, date := range dates {
					if *dataE2eDate == date {
						availableDatesMap[date] = struct{}{}
					}
				}
			}
		}
	}

	// Convert the map keys to a sorted slice
	availableDates := make([]string, 0, len(availableDatesMap))
	for date := range availableDatesMap {
		availableDates = append(availableDates, date)
	}
	sort.Strings(availableDates)

	// Return true if any dates are available, along with the list of unique available dates
	return len(availableDates) > 0, availableDates
}
