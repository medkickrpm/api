package worker

import "time"

func getMonthNumberFrom2023() int {
	year, month, _ := time.Now().UTC().Date()
	monthNumber := int(month)
	yearNumber := year

	return (yearNumber-2023)*12 + monthNumber
}
