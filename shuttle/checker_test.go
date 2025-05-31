package shuttle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWebChecker_CheckAvailabilityEndSeptLakeMoraine(t *testing.T) {
	// Arrange
	url := "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-09-29&endDate=2025-09-30&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-29T21:34:13.097&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D"
	checker := NewWebChecker()

	// Act
	result := checker.CheckAvailability(url)

	// Assert
	assert.True(t, result, "Expected available slots for URL: %s", url)
}

func TestWebChecker_CheckAvailabilityURLAug04_07_LakeLouise(t *testing.T) {
	// Arrange
	url := "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483089&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-04&endDate=2025-08-05&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-29T22:40:19.637&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D"
	checker := NewWebChecker()

	// Act
	result := checker.CheckAvailability(url)

	// Assert
	assert.False(t, result, "Expected available slots for URL: %s", url)
}

func TestWebChecker_CheckAvailabilityURL3(t *testing.T) {
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
	// Arrange
	url := "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-09-29&endDate=2025-09-30&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-29T21:34:13.097&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D"
	checker := NewWebChecker()

	// Act
	result := checker.CheckAvailability(url)
	// Assert
	assert.True(t, result, "Expected available slots for URL: %s", url)
}
