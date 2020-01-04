package rrule

func IsLeapYear(year int) int {
	if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		return 1
	}
	return 0
}
