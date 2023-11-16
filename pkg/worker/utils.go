package worker

import (
	"time"
)

func getMonthNumberFrom2023() int {
	year, month, _ := getStartTimeOfToday().Date()
	monthNumber := int(month)
	yearNumber := year

	return (yearNumber-2023)*12 + monthNumber
}

// Get start date time of a day in USA time zone
func getStartTimeOfToday() time.Time {
	loc, _ := time.LoadLocation("EST")
	nowTime := time.Now().UTC()
	year, month, day := nowTime.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

// Get End date time of a day in USA time zone
func getEndTimeOfToday() time.Time {
	loc, _ := time.LoadLocation("EST")
	nowTime := time.Now().UTC()
	year, month, day := nowTime.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, loc)
}

func getStartDateOfMonth() time.Time {
	loc, _ := time.LoadLocation("EST")
	nowTime := time.Now().UTC()
	year, month, _ := nowTime.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, loc)
}
