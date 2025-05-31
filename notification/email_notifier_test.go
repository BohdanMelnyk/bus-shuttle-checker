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
		"sender@example.com",
	)

	// Test sending notification
	id, err := notifier.SendNotification("https://test.com", "Test Location")
	
	// We expect an error since we're using dummy credentials
	assert.Error(t, err)
	assert.Empty(t, id)
}
