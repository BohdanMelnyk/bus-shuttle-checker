package shuttle

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type mockTransport struct {
	expectedURL     string
	expectedMethod  string
	expectedHeaders map[string]string
	expectedBody    string
	response        string
	t              *testing.T
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Verify URL
	if req.URL.String() != m.expectedURL {
		m.t.Errorf("Expected URL %s, got %s", m.expectedURL, req.URL.String())
	}

	// Verify method
	if req.Method != m.expectedMethod {
		m.t.Errorf("Expected method %s, got %s", m.expectedMethod, req.Method)
	}

	// Verify headers
	for key, expectedValue := range m.expectedHeaders {
		if actualValue := req.Header.Get(key); actualValue != expectedValue {
			m.t.Errorf("Expected header %s: %s, got %s", key, expectedValue, actualValue)
		}
	}

	// Verify body
	body, _ := io.ReadAll(req.Body)
	if string(body) != m.expectedBody {
		m.t.Errorf("Expected body %s, got %s", m.expectedBody, string(body))
	}

	// Return mock response
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(m.response)),
	}, nil
}

func TestHasAvailability(t *testing.T) {
	tests := []struct {
		name             string
		urlName          string
		resourceLocationID int
		startDate        string
		endDate          string
		resourceIDs      []int64
		bookingCategory  int
		expectedURL      string
		expectedHeaders  map[string]string
		expectedBody     string
		mockResponse     string
		wantAvailable    bool
		wantDates        []string
		wantErr          bool
	}{
		{
			name:             "Lake Morain Morning - October 11",
			urlName:          "Lake Morain Morning",
			resourceLocationID: -2147483642,
			startDate:        "2025-10-11",
			endDate:         "2025-10-11",
			resourceIDs:      []int64{-2147476652, -2147476634, -2147476641, -2147476655},
			bookingCategory:  9,
			expectedURL:      "https://reservation.pc.gc.ca/api/availability/dailyactivity?resourceLocationId=-2147483642&startDate=2025-10-11&endDate=2025-10-11&bookingCategoryId=9",
			expectedHeaders: map[string]string{
				"Accept":             "application/json",
				"Content-Type":       "application/json",
				"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
				"App-Language":       "en-CA",
				"app-version":        "5.98.197",
			},
			expectedBody: `[-2147476652,-2147476634,-2147476641,-2147476655]`,
			mockResponse: `[
				{"resourceId":-2147476652,"range":{"start":"2025-10-11T00:00:00Z","end":"2025-10-11T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476652,"startDate":"2025-10-11T00:00:00Z","endDate":"2025-10-12T00:00:00Z","remainingReservableQuota":2,"remainingTotalQuota":2,"closedQuota":0}},
				{"resourceId":-2147476634,"range":{"start":"2025-10-11T00:00:00Z","end":"2025-10-11T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476634,"startDate":"2025-10-11T00:00:00Z","endDate":"2025-10-12T00:00:00Z","remainingReservableQuota":2,"remainingTotalQuota":2,"closedQuota":0}},
				{"resourceId":-2147476641,"range":{"start":"2025-10-11T00:00:00Z","end":"2025-10-11T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476641,"startDate":"2025-10-11T00:00:00Z","endDate":"2025-10-12T00:00:00Z","remainingReservableQuota":2,"remainingTotalQuota":2,"closedQuota":0}},
				{"resourceId":-2147476655,"range":{"start":"2025-10-11T00:00:00Z","end":"2025-10-11T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476655,"startDate":"2025-10-11T00:00:00Z","endDate":"2025-10-12T00:00:00Z","remainingReservableQuota":2,"remainingTotalQuota":2,"closedQuota":0}}
			]`,
			wantAvailable: true,
			wantDates:     []string{"2025-10-11"},
			wantErr:       false,
		},
		{
			name:             "Lake Morain Midday - October 8",
			urlName:          "Lake Morain Midday",
			resourceLocationID: -2147483642,
			startDate:        "2025-10-08",
			endDate:         "2025-10-08",
			resourceIDs:      []int64{-2147476654, -2147476653, -2147476651, -2147476640, -2147471867, -2147471865, -2147471863, -2147471861},
			bookingCategory:  9,
			expectedURL:      "https://reservation.pc.gc.ca/api/availability/dailyactivity?resourceLocationId=-2147483642&startDate=2025-10-08&endDate=2025-10-08&bookingCategoryId=9",
			expectedHeaders: map[string]string{
				"Accept":             "application/json",
				"Content-Type":       "application/json",
				"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
				"App-Language":       "en-CA",
				"app-version":        "5.98.197",
			},
			expectedBody: `[-2147476654,-2147476653,-2147476651,-2147476640,-2147471867,-2147471865,-2147471863,-2147471861]`,
			mockResponse: `[
				{"resourceId":-2147476654,"range":{"start":"2025-10-08T00:00:00Z","end":"2025-10-08T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476654,"startDate":"2025-10-08T00:00:00Z","endDate":"2025-10-09T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147476653,"range":{"start":"2025-10-08T00:00:00Z","end":"2025-10-08T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476653,"startDate":"2025-10-08T00:00:00Z","endDate":"2025-10-09T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147476651,"range":{"start":"2025-10-08T00:00:00Z","end":"2025-10-08T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476651,"startDate":"2025-10-08T00:00:00Z","endDate":"2025-10-09T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147476640,"range":{"start":"2025-10-08T00:00:00Z","end":"2025-10-08T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476640,"startDate":"2025-10-08T00:00:00Z","endDate":"2025-10-09T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147471867,"range":{"start":"2025-10-08T00:00:00Z","end":"2025-10-08T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147471867,"startDate":"2025-10-08T00:00:00Z","endDate":"2025-10-09T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147471865,"range":{"start":"2025-10-08T00:00:00Z","end":"2025-10-08T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147471865,"startDate":"2025-10-08T00:00:00Z","endDate":"2025-10-09T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147471863,"range":{"start":"2025-10-08T00:00:00Z","end":"2025-10-08T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147471863,"startDate":"2025-10-08T00:00:00Z","endDate":"2025-10-09T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147471861,"range":{"start":"2025-10-08T00:00:00Z","end":"2025-10-08T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147471861,"startDate":"2025-10-08T00:00:00Z","endDate":"2025-10-09T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}}
			]`,
			wantAvailable: true,
			wantDates:     []string{"2025-10-08"},
			wantErr:       false,
		},
		{
			name:             "Lake Morain Morning",
			urlName:          "Lake Morain Morning",
			resourceLocationID: -2147483642,
			startDate:        "2025-08-05",
			endDate:         "2025-08-07",
			resourceIDs:      []int64{-2147476652, -2147476634, -2147476641, -2147476655},
			bookingCategory:  9,
			expectedURL:      "https://reservation.pc.gc.ca/api/availability/dailyactivity?resourceLocationId=-2147483642&startDate=2025-08-05&endDate=2025-08-07&bookingCategoryId=9",
			expectedHeaders: map[string]string{
				"Accept":             "application/json",
				"Content-Type":       "application/json",
				"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
				"App-Language":       "en-CA",
				"app-version":        "5.98.197",
			},
			expectedBody: `[-2147476652,-2147476634,-2147476641,-2147476655]`,
			mockResponse: `[
				{"resourceId":-2147476652,"range":{"start":"2025-08-05T00:00:00Z","end":"2025-08-05T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476652,"startDate":"2025-08-05T00:00:00Z","endDate":"2025-08-06T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147476634,"range":{"start":"2025-08-05T00:00:00Z","end":"2025-08-05T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476634,"startDate":"2025-08-05T00:00:00Z","endDate":"2025-08-06T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147476641,"range":{"start":"2025-08-05T00:00:00Z","end":"2025-08-05T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476641,"startDate":"2025-08-05T00:00:00Z","endDate":"2025-08-06T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147476655,"range":{"start":"2025-08-05T00:00:00Z","end":"2025-08-05T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476655,"startDate":"2025-08-05T00:00:00Z","endDate":"2025-08-06T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}}
			]`,
			wantAvailable: true,
			wantDates:     []string{"2025-08-05"},
			wantErr:       false,
		},
		{
			name:             "Lake Morain Midday",
			urlName:          "Lake Morain Midday",
			resourceLocationID: -2147483642,
			startDate:        "2025-08-05",
			endDate:         "2025-08-07",
			resourceIDs:      []int64{-2147476651, -2147476653},
			bookingCategory:  9,
			expectedURL:      "https://reservation.pc.gc.ca/api/availability/dailyactivity?resourceLocationId=-2147483642&startDate=2025-08-05&endDate=2025-08-07&bookingCategoryId=9",
			expectedHeaders: map[string]string{
				"Accept":             "application/json",
				"Content-Type":       "application/json",
				"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
				"App-Language":       "en-CA",
				"app-version":        "5.98.197",
			},
			expectedBody: `[-2147476651,-2147476653]`,
			mockResponse: `[
				{"resourceId":-2147476651,"range":{"start":"2025-08-05T00:00:00Z","end":"2025-08-05T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476651,"startDate":"2025-08-05T00:00:00Z","endDate":"2025-08-06T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147476653,"range":{"start":"2025-08-05T00:00:00Z","end":"2025-08-05T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147476653,"startDate":"2025-08-05T00:00:00Z","endDate":"2025-08-06T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}}
			]`,
			wantAvailable: true,
			wantDates:     []string{"2025-08-05"},
			wantErr:       false,
		},
		{
			name:             "Lake O'Hara",
			urlName:          "Lake O'Hara",
			resourceLocationID: -2147483536,
			startDate:        "2025-08-05",
			endDate:         "2025-08-07",
			resourceIDs:      []int64{-2147479230, -2147479229},
			bookingCategory:  10,
			expectedURL:      "https://reservation.pc.gc.ca/api/availability/dailyactivity?resourceLocationId=-2147483536&startDate=2025-08-05&endDate=2025-08-07&bookingCategoryId=10",
			expectedHeaders: map[string]string{
				"Accept":             "application/json",
				"Content-Type":       "application/json",
				"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
				"App-Language":       "en-CA",
				"app-version":        "5.98.197",
			},
			expectedBody: `[-2147479230,-2147479229]`,
			mockResponse: `[
				{"resourceId":-2147479230,"range":{"start":"2025-08-05T00:00:00Z","end":"2025-08-05T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147479230,"startDate":"2025-08-05T00:00:00Z","endDate":"2025-08-06T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}},
				{"resourceId":-2147479229,"range":{"start":"2025-08-05T00:00:00Z","end":"2025-08-05T00:00:00Z"},"availabilityResult":{"resultCode":0,"resourceId":-2147479229,"startDate":"2025-08-05T00:00:00Z","endDate":"2025-08-06T00:00:00Z","remainingReservableQuota":1,"remainingTotalQuota":1,"closedQuota":0}}
			]`,
			wantAvailable: true,
			wantDates:     []string{"2025-08-05"},
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock transport
			mockTransport := &mockTransport{
				expectedURL:     tt.expectedURL,
				expectedMethod:  "POST",
				expectedHeaders: tt.expectedHeaders,
				expectedBody:    tt.expectedBody,
				response:       tt.mockResponse,
				t:              t,
			}

			// Create client with mock transport
			client := &APIClient{
				client: &http.Client{
					Transport: mockTransport,
				},
			}

			// Call HasAvailability
			gotAvailable, gotDates, err := client.HasAvailability(
				tt.urlName,
				tt.resourceLocationID,
				tt.startDate,
				tt.endDate,
				tt.resourceIDs,
				tt.bookingCategory,
			)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("HasAvailability() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check availability
			if gotAvailable != tt.wantAvailable {
				t.Errorf("HasAvailability() available = %v, want %v", gotAvailable, tt.wantAvailable)
			}

			// Check dates
			if !reflect.DeepEqual(gotDates, tt.wantDates) {
				t.Errorf("HasAvailability() dates = %v, want %v", gotDates, tt.wantDates)
			}
		})
	}
} 