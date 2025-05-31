package notification

// Notifier is an interface for sending notifications
type Notifier interface {
	// SendNotification sends a notification about availability
	// parameters:
	//   - url: the URL where availability was found
	//   - locationName: the name of the location that has availability
	// returns:
	//   - id: a unique identifier for the sent notification (if applicable)
	//   - error: any error that occurred during notification sending
	SendNotification(url string, locationName string) (id string, err error)
}
