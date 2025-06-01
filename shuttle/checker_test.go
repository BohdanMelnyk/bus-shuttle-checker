package shuttle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWebChecker_CheckAvailabilityEndSeptLakeMoraine(t *testing.T) {
	t.Skip("Skipping legacy availability check to avoid external API calls. Remove t.Skip() to run this test manually.")
	// Arrange
	url := "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-09-29&endDate=2025-09-30&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-29T21:34:13.097&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D"
	checker := NewWebChecker()

	// Act
	result := checker.CheckAvailability(url)

	// Assert
	assert.True(t, result, "Expected available slots for URL: %s", url)
}

func TestWebChecker_CheckAvailabilityURLAug04_07_LakeLouise(t *testing.T) {
	t.Skip("Skipping legacy availability check to avoid external API calls. Remove t.Skip() to run this test manually.")
	// Arrange
	url := "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483089&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-04&endDate=2025-08-05&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-29T22:40:19.637&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D"
	checker := NewWebChecker()

	// Act
	result := checker.CheckAvailability(url)

	// Assert
	assert.False(t, result, "Expected available slots for URL: %s", url)
}

func TestWebChecker_CheckAvailabilityURL3(t *testing.T) {
	t.Skip("Skipping legacy availability check to avoid external API calls. Remove t.Skip() to run this test manually.")
	// Arrange
	url := "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-04&endDate=2025-08-05&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-29T23:02:31.428&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D"
	checker := NewWebChecker()

	// Act
	result := checker.CheckAvailability(url)

	// Assert
	assert.False(t, result, "Expected not-available slots for URL: %s", url)
}

// Also test the compatibility function
func TestLegacyCheckAvailability(t *testing.T) {
	t.Skip("Skipping legacy availability check to avoid external API calls. Remove t.Skip() to run this test manually.")
	// Arrange
	url := "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-09-29&endDate=2025-09-30&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-29T21:34:13.097&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D"
	checker := NewWebChecker()

	// Act
	result := checker.CheckAvailability(url)
	// Assert
	assert.True(t, result, "Expected available slots for URL: %s", url)
}

func TestWebChecker_CheckAvailabilityForDates_Available(t *testing.T) {
	// Arrange
	url := "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-10-07&endDate=2025-10-08&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,1,null%5D%5D&searchTime=2025-06-01T14:02:43.225&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D"
	checker := NewWebChecker()
	dates := []string{"2025-10-12"}

	// Act
	available, availableDates := checker.CheckAvailabilityForDates(url, dates)

	// Assert
	assert.True(t, available, "Expected available slots for URL: %s and dates: %v", url, dates)
	assert.Equal(t, []string{"2025-10-12"}, availableDates, "Expected available dates to match")
}

func TestWebChecker_CheckAvailabilityForDates_NotAvailable(t *testing.T) {
	// Arrange
	url := "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-04&endDate=2025-08-05&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,1,null%5D%5D&searchTime=2025-06-01T14:03:40.358&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D"
	checker := NewWebChecker()
	dates := []string{"2025-06-08", "2025-07-08"}

	// Act
	available, availableDates := checker.CheckAvailabilityForDates(url, dates)

	// Assert
	assert.False(t, available, "Expected no available slots for URL: %s and dates: %v", url, dates)
	assert.Empty(t, availableDates, "Expected no available dates")
}
