package rrule

import (
	"testing"
	"time"
)

// https://tools.ietf.org/html/rfc5545#page-123
func Test_rfc_5545_page_123_count(t *testing.T) {
	// ==> (1997 9:00 AM EDT) September 2-11
	var value = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=DAILY;COUNT=10"

	rule, err := AssertToStringMatchesInput(value)
	if err != nil {
		t.Log("Failed to parse", err)
		t.Fail()
	}

	if rule.Frequency != DAILY {
		t.Log("Failed to parse frequency")
		t.Fail()
	}

	if rule.Count != 10 {
		t.Log("Failed to parse count correctly.")
		t.Fail()
	}
}

// https://tools.ietf.org/html/rfc5545#page-123
func Test_rfc_5545_page_123_until(t *testing.T) {
	// ==> (1997 9:00 AM EDT) September 2-30;October 1-25
	//     (1997 9:00 AM EST) October 26-31;November 1-30;December 1-23
	var value = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=DAILY;UNTIL=19971224T000000Z"

	rule, err := AssertToStringMatchesInput(value)
	if err != nil {
		t.Log("Failed to parse", err)
		t.Fail()
	}

	if rule.Frequency != DAILY {
		t.Log("Failed to parse frequency")
		t.Fail()
	}

	targetUntil := time.Date(1997, 12, 24, 0, 0, 0, 0, time.UTC)

	if !rule.Until.Equal(targetUntil) {
		t.Log("Failed to parse until correctly.")
		t.Fail()
	}
}

// https://tools.ietf.org/html/rfc5545#page-129
func Test_EveryFridayThe13thForever(t *testing.T) {
	var value = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"EXDATE;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=MONTHLY;BYDAY=FR;BYMONTHDAY=13"

	rule, err := AssertToStringMatchesInput(value)
	if err != nil {
		t.Log("Failed to parse", err)
		t.Fail()
	}

	targetLocation, _ := time.LoadLocation("America/New_York")
	if rule.DtStart.Location().String() != targetLocation.String() {
		t.Log("Locations don't match")
		t.Fail()
	}

	if len(rule.ExceptionsToRule) != 1 {
		t.Log("Failed to parse DTSTART")
		t.Fail()
	}

	if !rule.ExceptionsToRule[0].Equal(
		time.Date(1997, 9, 2, 9, 0, 0, 0, targetLocation)) {
		t.Log("Failed to parse EXDATE")
		t.Fail()
	}

	if rule.Frequency != MONTHLY {
		t.Log("Failed to parse frequency")
		t.Fail()
	}

}

// https://tools.ietf.org/html/rfc5545#page-45
func Test_Interval_and_by_many(t *testing.T) {
	var value = "DTSTART;TZID=America/New_York:19970105T083000\n" +
		"RRULE:FREQ=YEARLY;INTERVAL=2;BYMONTH=1;BYDAY=SU;" +
		"BYHOUR=8,9;BYMINUTE=30"

	rule, err := AssertToStringMatchesInput(value)
	if err != nil {
		t.Log("Failed to parse rule")
		t.Fail()
	}

	if rule.Interval != 2 {
		t.Log("Failed to parse Interval")
		t.Fail()
	}

	if rule.ByDay[0].Weekday != time.Sunday || rule.ByDay[0].Offset != 0 {
		t.Log("Didnt parse ByDay correctly.")
		t.Fail()
	}

	if rule.ByHour[0] != 8 || rule.ByHour[1] != 9 {
		t.Log("Failed to parse ByHour")
		t.Fail()
	}

	if rule.ByMonth[0] != 1 {
		t.Log("Failed to parse ByMonth")
		t.Fail()
	}

	if rule.ByMinute[0] != 30 {
		t.Log("Failed to parse byminute")
		t.Fail()
	}

}

func Test_Parse_and_back_to_string(t *testing.T) {
	var value = "RRULE:FREQ=WEEKLY;WKST=MO;BYDAY=MO,TU,WE,TH,FR"
	_, err := AssertToStringMatchesInput(value)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

}

func Test_Parse_set_until_render(t *testing.T) {
	var value = "RRULE:FREQ=WEEKLY;WKST=MO;BYDAY=MO,TU,WE,TH,FR"
	rule, err := AssertToStringMatchesInput(value)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	rule.Until = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)

	rule2, err := Parse(rule.String())

	if !rule2.Until.Equal(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Log("Failed to set new until value")
		t.Fail()
	}
}

func Test_Parse_set_until_from_forever(t *testing.T) {
	var value = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=DAILY"
	rule, err := AssertToStringMatchesInput(value)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	targetLocation, _ := time.LoadLocation("America/New_York")
	endDate := time.Date(1997, time.December, 1, 0, 0, 0, 0, targetLocation)

	rule.Until = endDate

	var iter *RecurrenceIterator = rule.Iterator()
	var event time.Time

	for iter.Step(&event) {
		if event.After(endDate) {
			t.Log("Event created after set Until")
			t.Fail()
		}
	}

	if iter.ReturnCounter <= 0 {
		t.Log("Iterator never returned any events")
		t.Fail()
	}
}

func Test_Parse_PositiveIndicatorIgnored(t *testing.T) {
	var withPos = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=MONTHLY;BYDAY=+1MO,+1TU"
	wpRule, err := AssertToStringMatchesInput(withPos)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	var noPos = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=MONTHLY;BYDAY=1MO,1TU"
	npRule, err := AssertToStringMatchesInput(noPos)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if wpRule.String() != npRule.String() {
		t.Log("positive('+') indicator not ignored")
		t.Fail()
	}
}
