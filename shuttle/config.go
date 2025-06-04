package shuttle

type Location struct {
	Name             string
	LocationID       int
	ResourceIDs      []int64
	Dates           []string
	BookingCategory int // 9 for Moraine Lake, 10 for Lake O'Hara
}

var (
	LakeMorainMorning = Location{
		Name:             "Lake Morain Morning",
		LocationID:       -2147483642,
		ResourceIDs:      []int64{-2147476652, -2147476634, -2147476641, -2147476655},
		Dates:           []string{"2025-08-05", "2025-08-06", "2025-08-07"},
		BookingCategory: 9,
	}

	LakeMorainMidday = Location{
		Name:             "Lake Morain Midday",
		LocationID:       -2147483642,
		ResourceIDs:      []int64{-2147476651, -2147476653},
		Dates:           []string{"2025-08-05", "2025-08-06", "2025-08-07"},
		BookingCategory: 9,
	}

	LakeOHara = Location{
		Name:             "Lake O'Hara",
		LocationID:       -2147483536,
		ResourceIDs:      []int64{-2147479230, -2147479229},
		Dates:           []string{"2025-08-05", "2025-08-06", "2025-08-07"},
		BookingCategory: 10,
	}

	// List of all locations to check
	Locations = []Location{
		LakeMorainMorning,
		LakeMorainMidday,
		LakeOHara,
	}
) 