package shuttle

// AvailabilityChecker defines the interface for checking shuttle availability
type AvailabilityChecker interface {
	// CheckAvailability checks if a shuttle slot is available for the given URL
	// It returns true if a slot is available, false otherwise
	CheckAvailability(url string) bool

	// CheckAvailabilityForDates checks if a shuttle slot is available for the given URL and specific dates
	// It returns true if a slot is available on any of the specified dates, false otherwise
	// Additionally, it returns a list of dates for which slots are available
	CheckAvailabilityForDates(url string, dates []string) (bool, []string)
}
