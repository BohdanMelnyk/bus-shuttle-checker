package shuttle

// ShuttleInfo contains information about a shuttle reservation including URL and dates of interest
type ShuttleInfo struct {
	URL   string
	Dates []string // Dates in format "YYYY-MM-DD" that we are monitoring
}

// ShuttleURLs maps descriptive names to shuttle reservation information
var ShuttleURLs = map[string]ShuttleInfo{
	"MoraineLakeAlpineStart": {
		URL:   "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483092&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T00:13:45.837&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
		Dates: []string{"2025-08-06", "2025-08-07", "2025-08-08"},
	},
	"MoraineLakeMorning": {
		URL:   "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T00:13:45.837&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
		Dates: []string{"2025-08-06", "2025-08-07", "2025-08-08"},
	},
	"MoraineLakeMidday": {
		URL:   "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483087&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T00:13:45.837&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
		Dates: []string{"2025-08-06", "2025-08-07", "2025-08-08"},
	},

	"LouiseLakeAlpineStart": {
		URL:   "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483091&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T13:37:19.879&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
		Dates: []string{"2025-08-06", "2025-08-07", "2025-08-08"},
	},
	"LouiseLakeMorning": {
		URL:   "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483089&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T13:37:19.879&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
		Dates: []string{"2025-08-06", "2025-08-07", "2025-08-08"},
	},
	"LouiseLakeMidday": {
		URL:   "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483086&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-04&endDate=2025-08-05&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T13:38:42.780&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
		Dates: []string{"2025-08-06", "2025-08-07", "2025-08-08"},
	},

	"O'HaraLake": {
		URL:   "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483328&searchTabGroupId=3&bookingCategoryId=10&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T00:14:32.599&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483536&filterData=%7B%7D",
		Dates: []string{"2025-08-06", "2025-08-07", "2025-08-08"},
	},
}

// GetAllURLs returns all shuttle URLs as a map of name to URL
func GetAllURLs() map[string]string {
	result := make(map[string]string)
	for name, info := range ShuttleURLs {
		result[name] = info.URL
	}
	return result
}

// GetURL returns the URL for the given shuttle name
func GetURL(name string) (string, bool) {
	info, exists := ShuttleURLs[name]
	if exists {
		return info.URL, true
	}
	return "", false
}

// GetDates returns the dates of interest for the given shuttle name
func GetDates(name string) ([]string, bool) {
	info, exists := ShuttleURLs[name]
	if exists {
		return info.Dates, true
	}
	return nil, false
}

// GetShuttleInfo returns the complete shuttle information for the given name
func GetShuttleInfo(name string) (ShuttleInfo, bool) {
	info, exists := ShuttleURLs[name]
	return info, exists
}

// GetAllShuttleInfo returns all shuttle information
func GetAllShuttleInfo() map[string]ShuttleInfo {
	return ShuttleURLs
}
