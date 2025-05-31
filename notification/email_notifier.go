package notification

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"time"
)

// EmailNotifier implements the Notifier interface using Mailgun for email delivery
type EmailNotifier struct {
	Domain    string
	APIKey    string
	Recipient string
	Sender    string
}

// NewEmailNotifier creates a new EmailNotifier with the given configuration
func NewEmailNotifier(domain, apiKey, recipient, sender string) *EmailNotifier {
	return &EmailNotifier{
		Domain:    domain,
		APIKey:    apiKey,
		Recipient: recipient,
		Sender:    sender,
	}
}

// SendNotification sends an email notification about an available shuttle slot
func (e *EmailNotifier) SendNotification(url string, locationName string) (string, error) {
	mg := mailgun.NewMailgun(e.Domain, e.APIKey)
	m := mailgun.NewMessage(
		e.Sender,
		fmt.Sprintf("Shuttle Slot Available for %s", locationName),
		fmt.Sprintf("A shuttle slot is now available for booking at %s.\n\nBooking URL: %s", locationName, url),
		e.Recipient,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}
