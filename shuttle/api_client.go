package shuttle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type AvailabilityRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type AvailabilityResult struct {
	ResultCode              int       `json:"resultCode"`
	ResourceID             int       `json:"resourceId"`
	StartDate              string    `json:"startDate"`
	EndDate                string    `json:"endDate"`
	RemainingReservableQuota float64  `json:"remainingReservableQuota"`
	RemainingTotalQuota     float64  `json:"remainingTotalQuota"`
	ClosedQuota            float64  `json:"closedQuota"`
	Exception              any      `json:"exception"`
}

type ResourceAvailability struct {
	ResourceID        int               `json:"resourceId"`
	Range            AvailabilityRange `json:"range"`
	AvailabilityResult AvailabilityResult `json:"availabilityResult"`
}

type APIClient struct {
	client *http.Client
}

func NewAPIClient() *APIClient {
	return &APIClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *APIClient) CheckAvailability(locationID int, startDate, endDate string, resourceIDs []int) ([]ResourceAvailability, error) {
	url := fmt.Sprintf("https://reservation.pc.gc.ca/api/availability/dailyactivity?resourceLocationId=%d&startDate=%s&endDate=%s&bookingCategoryId=10",
		locationID, startDate, endDate)

	// Convert resource IDs to JSON
	resourceIDsJSON, err := json.Marshal(resourceIDs)
	if err != nil {
		return nil, fmt.Errorf("error marshaling resource IDs: %w", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(resourceIDsJSON)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add required headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	req.Header.Add("App-Language", "en-CA")
	req.Header.Add("app-version", "5.98.197")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var availabilities []ResourceAvailability
	if err := json.Unmarshal(body, &availabilities); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return availabilities, nil
}

// HasAvailability checks if specific resources have available quota
func (c *APIClient) HasAvailability(urlName string, resourceLocationId int, startDate, endDate string, resourceIds []int64, bookingCategory int) (bool, []string, error) {
	url := fmt.Sprintf("https://reservation.pc.gc.ca/api/availability/dailyactivity?resourceLocationId=%d&startDate=%s&endDate=%s&bookingCategoryId=%d",
		resourceLocationId, startDate, endDate, bookingCategory)

	// Convert resource IDs to JSON
	resourceIDsJSON, err := json.Marshal(resourceIds)
	if err != nil {
		return false, nil, fmt.Errorf("error marshaling resource IDs: %w", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(resourceIDsJSON)))
	if err != nil {
		return false, nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add required headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	req.Header.Add("App-Language", "en-CA")
	req.Header.Add("app-version", "5.98.197")

	resp, err := c.client.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// First try to parse with the standard time format
	var availabilities []ResourceAvailability
	if err := json.Unmarshal(body, &availabilities); err != nil {
		// If that fails, try with a custom time format
		var rawData []map[string]json.RawMessage
		if err := json.Unmarshal(body, &rawData); err != nil {
			return false, nil, fmt.Errorf("error unmarshaling response: %w", err)
		}

		// Process each item manually
		availabilities = make([]ResourceAvailability, len(rawData))
		for i, raw := range rawData {
			var resourceID int
			if err := json.Unmarshal(raw["resourceId"], &resourceID); err != nil {
				return false, nil, fmt.Errorf("error unmarshaling resourceId: %w", err)
			}

			var rangeData map[string]string
			if err := json.Unmarshal(raw["range"], &rangeData); err != nil {
				return false, nil, fmt.Errorf("error unmarshaling range: %w", err)
			}

			start, err := time.Parse("2006-01-02T15:04:05", rangeData["start"])
			if err != nil {
				return false, nil, fmt.Errorf("error parsing start time: %w", err)
			}

			end, err := time.Parse("2006-01-02T15:04:05", rangeData["end"])
			if err != nil {
				return false, nil, fmt.Errorf("error parsing end time: %w", err)
			}

			var result AvailabilityResult
			if err := json.Unmarshal(raw["availabilityResult"], &result); err != nil {
				return false, nil, fmt.Errorf("error unmarshaling availabilityResult: %w", err)
			}

			availabilities[i] = ResourceAvailability{
				ResourceID: resourceID,
				Range: AvailabilityRange{
					Start: start,
					End:   end,
				},
				AvailabilityResult: result,
			}
		}
	}

	availableDates := make(map[string]struct{})
	
	// Group availabilities by date to ensure both resources are available
	dateResources := make(map[string]map[int]float64)
	
	for _, avail := range availabilities {
		date := avail.Range.Start.Format("2006-01-02")
		if _, exists := dateResources[date]; !exists {
			dateResources[date] = make(map[int]float64)
		}
		dateResources[date][avail.ResourceID] = avail.AvailabilityResult.RemainingReservableQuota
	}

	// Check each date for all required resource IDs
	for date, resources := range dateResources {
		allResourcesAvailable := true
		for _, resourceID := range resourceIds {
			quota, exists := resources[int(resourceID)]
			if !exists || quota <= 0 {
				allResourcesAvailable = false
				break
			}
		}
		if allResourcesAvailable {
			availableDates[date] = struct{}{}
		}
	}

	// Convert map to sorted slice
	dates := make([]string, 0, len(availableDates))
	for date := range availableDates {
		dates = append(dates, date)
	}

	return len(dates) > 0, dates, nil
} 