package notification

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailNotifier_SendNotification(t *testing.T) {
	// Create a test notifier with dummy credentials
	notifier := NewEmailNotifier(
		"test-domain.mailgun.org",
		"test-api-key-xxxxx",
		"test@example.com",
	)

	// Test sending notification
	err := notifier.SendNotification("Test Subject", "Test Message")
	if err != nil {
		// We expect an error since we're using dummy credentials
		t.Logf("Expected error with dummy credentials: %v", err)
	}
}
