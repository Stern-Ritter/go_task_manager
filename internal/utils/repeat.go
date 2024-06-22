package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/Stern-Ritter/go_task_manager/internal/errors"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	isRepeatValid, err := ValidateRepeat(repeat)
	if err != nil || !isRepeatValid {
		return "", errors.NewInvalidRepeatFormat("invalid task repeat format", err)
	}

	d, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.NewInvalidDateFormat("invalid task date format", err)
	}

	parts := parseRepeat(repeat)
	switch parts["type"] {
	case "y":
		return nextY(now, d)
	case "d":
		return nextD(now, d, parts["value"])
	case "w":
		return nextW(now, d, parts["value"])
	case "m":
		return nextM(now, d, parts["value"])
	default:
		return "", errors.NewInvalidRepeatFormat("invalid task repeat format", nil)
	}
}

func nextY(now time.Time, date time.Time) (string, error) {
	res := date.AddDate(1, 0, 0)
	for res.Before(now) {
		res = res.AddDate(1, 0, 0)
	}

	return res.Format("20060102"), nil
}

func nextD(now time.Time, date time.Time, value string) (string, error) {
	daysCount, err := strconv.Atoi(value)
	if err != nil {
		return "", err
	}

	res := date
	res = res.AddDate(0, 0, daysCount)
	for res.Before(now) {
		res = res.AddDate(0, 0, daysCount)
	}

	return res.Format("20060102"), nil
}

func nextW(now time.Time, date time.Time, value string) (string, error) {
	res := getMaxDate(now, date)
	weekDays, err := parseWeekDaysValue(value)
	if err != nil {
		return "", err
	}

	res = res.AddDate(0, 0, 1)
	for !contains(weekDays, parseWeekDay(res.Weekday())) {
		res = res.AddDate(0, 0, 1)
	}

	return res.Format("20060102"), nil
}

func nextM(now time.Time, date time.Time, value string) (string, error) {
	res := getMaxDate(now, date)
	days, months, err := parseMonthsDaysValue(value)
	if err != nil {
		return "", err
	}

	res = res.AddDate(0, 0, 1)
	for !(checkMonthDays(days, res.Day(), daysIn(res.Month(), res.Year())) &&
		(len(months) == 0 || contains(months, int(res.Month())))) {
		res = res.AddDate(0, 0, 1)
	}
	return res.Format("20060102"), nil
}

func parseRepeat(repeat string) map[string]string {
	result := make(map[string]string)
	parts := strings.Split(repeat, " ")
	result["type"] = parts[0]
	result["value"] = strings.Join(parts[1:], " ")
	return result
}

func parseWeekDaysValue(value string) ([]int, error) {
	parts := strings.Split(value, ",")
	res := make([]int, len(parts))
	for idx, el := range parts {
		num, err := strconv.Atoi(el)
		if err != nil {
			return []int{}, err
		}
		res[idx] = num
	}
	return res, nil
}

func parseMonthsDaysValue(value string) ([]int, []int, error) {
	parts := strings.Split(value, " ")
	daysPats := strings.Split(parts[0], ",")
	monthsParts := []string{}
	if len(parts) > 1 {
		monthsParts = strings.Split(parts[1], ",")
	}

	days := make([]int, len(daysPats))
	for idx, el := range daysPats {
		num, err := strconv.Atoi(el)
		if err != nil {
			return []int{}, []int{}, err
		}
		days[idx] = num
	}

	months := make([]int, len(monthsParts))
	for idx, el := range monthsParts {
		num, err := strconv.Atoi(el)
		if err != nil {
			return []int{}, []int{}, err
		}
		months[idx] = num
	}

	return days, months, nil
}

func checkMonthDays(days []int, currentDay int, monthDaysCount int) bool {
	for _, day := range days {
		if day > 0 && day == currentDay {
			return true
		} else if (monthDaysCount + day + 1) == currentDay {
			return true
		}
	}
	return false
}
