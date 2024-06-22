package utils

import "regexp"

const (
	RepeatPattern = "(^y$)|" +
		"(^d\\s([0123]?[0-9]?[0-9]?|400)$)|" +
		"(^w\\s[1-7]{1}(,[1-7]){0,6}$)|" +
		"(^m\\s([012]?[0-9]?|3[01]|-[12]{1}){1}(,([012]?[0-9]?|3[01]|-[12]{1})){0,30}(\\s(([0]?[0-9])|1[012]){1}(,(([0]?[0-9])|1[012])){0,11})?$)"
	SearchDatePatter = "(0[1-9]|[12][0-9]|3[01])\\.(0[1-9]|1[1,2])\\.(19|20)\\d{2}"
)

func ValidateRepeat(repeat string) (bool, error) {
	return regexp.MatchString(RepeatPattern, repeat)
}

func ValidateSearchDate(searchDate string) (bool, error) {
	return regexp.MatchString(SearchDatePatter, searchDate)
}
