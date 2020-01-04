package rrule

import (
	"log"
	"testing"
)

func Test_leap_year_negative(t *testing.T) {
	var not_leap_years []int = []int{1700, 1800, 1900, 2100, 2200, 2300, 2500, 2600}
	for _, year := range not_leap_years {
		if IsLeapYear(year) != 0 {
			log.Printf("%d isn't leap year.\n", year)
			t.Fail()
		}
	}
}

func Test_leap_year_positive(t *testing.T) {
	var leap_years []int = []int{1600, 2000, 2400}
	for _, year := range leap_years {
		if IsLeapYear(year) != 1 {
			log.Printf("%d is a leap year.\n", year)
			t.Fail()
		}
	}
}
