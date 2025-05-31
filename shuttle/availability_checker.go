package shuttle

// AvailabilityChecker defines the interface for checking shuttle availability
type AvailabilityChecker interface {
	// CheckAvailability checks if a shuttle slot is available for the given URL
	// It returns true if a slot is available, false otherwise
	CheckAvailability(url string) bool
}
