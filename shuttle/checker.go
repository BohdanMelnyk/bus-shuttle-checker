package shuttle

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// WebChecker implements the AvailabilityChecker interface using web scraping
type WebChecker struct{}

// NewWebChecker creates a new WebChecker instance
func NewWebChecker() *WebChecker {
	return &WebChecker{}
}

// getBrowserPath returns the appropriate browser path based on the operating system
func getBrowserPath() string {
	// First check environment variable
	if path := os.Getenv("CHROME_BIN"); path != "" {
		return path
	}

	// Then check common locations based on OS
	switch runtime.GOOS {
	case "darwin":
		paths := []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/usr/bin/chromium",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	case "linux":
		paths := []string{
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	case "windows":
		paths := []string{
			os.Getenv("PROGRAMFILES") + "\\Google\\Chrome\\Application\\chrome.exe",
			os.Getenv("PROGRAMFILES(X86)") + "\\Google\\Chrome\\Application\\chrome.exe",
			os.Getenv("LOCALAPPDATA") + "\\Google\\Chrome\\Application\\chrome.exe",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	// If no browser is found, let Rod handle it
	return ""
}

// CheckAvailability checks if a shuttle slot is available for the given URL
// It returns true if a slot is available, false otherwise
func (w *WebChecker) CheckAvailability(url string) bool {
	// Get browser path
	browserPath := getBrowserPath()
	log.Printf("Using browser at: %s", browserPath)

	// Configure launcher
	l := launcher.New()
	if browserPath != "" {
		l = l.Bin(browserPath)
	}

	// Set common flags for stability
	l = l.
		Set("no-sandbox").
		Set("disable-dev-shm-usage").
		Set("disable-gpu").
		Set("disable-software-rasterizer").
		Set("disable-dbus").
		Set("disable-setuid-sandbox").
		Set("single-process").
		Set("no-zygote").
		Set("no-first-run").
		Set("disable-extensions").
		Set("disable-background-networking").
		Set("disable-background-timer-throttling").
		Set("disable-backgrounding-occluded-windows").
		Set("disable-breakpad").
		Set("disable-component-extensions-with-background-pages").
		Set("disable-features", "TranslateUI,BlinkGenPropertyTrees").
		Set("disable-ipc-flooding-protection").
		Set("enable-automation").
		Headless(true)

	// Add retry logic for browser launch
	var browser *rod.Browser
	var err error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to launch browser (attempt %d/%d)...", i+1, maxRetries)
		controlURL := l.MustLaunch()
		browser = rod.New().
			ControlURL(controlURL).
			Timeout(30 * time.Second)

		err = browser.Connect()
		if err == nil {
			log.Printf("Successfully connected to browser")
			break
		}
		log.Printf("Failed to connect to browser (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		log.Printf("Failed to connect to browser after %d attempts: %v", maxRetries, err)
		return false
	}
	defer func() {
		if err := browser.Close(); err != nil {
			log.Printf("Error closing browser: %v", err)
		}
	}()

	// Create page with timeout
	log.Printf("Creating new page...")
	page := browser.MustPage()
	defer func() {
		if err := page.Close(); err != nil {
			log.Printf("Error closing page: %v", err)
		}
	}()

	// Navigate to URL with error handling
	log.Printf("Navigating to URL: %s", url)
	if err := page.Navigate(url); err != nil {
		log.Printf("Error navigating to URL: %v", err)
		return false
	}

	// Wait for the page to load completely
	log.Printf("Waiting for page to load...")
	if err := page.WaitLoad(); err != nil {
		log.Printf("Error waiting for page to load: %v", err)
		return false
	}

	// Additional wait to ensure all dynamic content is loaded
	log.Printf("Waiting for page to become idle...")
	if err := page.WaitIdle(30 * time.Second); err != nil {
		log.Printf("Error waiting for page idle: %v", err)
		return false
	}

	// Wait for the specific element indicating the page is fully loaded
	log.Printf("Waiting for chart-cell element...")
	element, err := page.Element(".chart-cell")
	if err != nil {
		log.Printf("Error finding chart-cell element: %v", err)
		return false
	}

	if err := element.WaitVisible(); err != nil {
		log.Printf("Error waiting for element visibility: %v", err)
		return false
	}

	// Extract the page source
	pageSource := page.MustHTML()

	// Check for availability based on aria-label containing "Available"
	available := strings.Contains(pageSource, "Departures Available")
	log.Printf("Availability check result: %v", available)
	return available
}
