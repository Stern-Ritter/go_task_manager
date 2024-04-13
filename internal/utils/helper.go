package utils

import "time"

func contains(arr []int, value int) bool {
	for _, el := range arr {
		if el == value {
			return true
		}
	}
	return false
}

func getMaxDate(dates ...time.Time) time.Time {
	max := dates[0]
	for i := 1; i < len(dates); i++ {
		if max.Before(dates[i]) {
			max = dates[i]
		}
	}
	return max
}

func parseWeekDay(weekDay time.Weekday) int {
	if weekDay == time.Sunday {
		return 7
	} else {
		return int(weekDay)
	}
}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
