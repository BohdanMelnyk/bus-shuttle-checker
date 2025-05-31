package shuttle

// ShuttleURLs maps descriptive names to shuttle reservation URLs
var ShuttleURLs = map[string]string{
	"MoraineLakeAlpineStart": "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483092&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T00:13:45.837&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
	"MoraineLakeMorning":     "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483090&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T00:13:45.837&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
	"MoraineLakeMidday":      "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483087&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T00:13:45.837&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",

	"LouiseLakeAlpineStart": "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483091&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T13:37:19.879&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
	"LouiseLakeMorning":     "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483089&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T13:37:19.879&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",
	"LouiseLakeMidday":      "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483086&searchTabGroupId=3&bookingCategoryId=9&startDate=2025-08-04&endDate=2025-08-05&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T13:38:42.780&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483642&filterData=%7B%7D",

	"O'HaraLake": "https://reservation.pc.gc.ca/create-booking/results?mapId=-2147483328&searchTabGroupId=3&bookingCategoryId=10&startDate=2025-08-05&endDate=2025-08-06&nights=1&isReserving=true&peopleCapacityCategoryCounts=%5B%5B-32767,null,2,null%5D%5D&searchTime=2025-05-30T00:14:32.599&groupHoldUid=&flexibleSearch=%5Bfalse,false,null,1%5D&resourceLocationId=-2147483536&filterData=%7B%7D",
}

// GetAllURLs returns all shuttle URLs as a map of name to URL
func GetAllURLs() map[string]string {
	return ShuttleURLs
}

// GetURL returns the URL for the given shuttle name
func GetURL(name string) (string, bool) {
	url, exists := ShuttleURLs[name]
	return url, exists
}
