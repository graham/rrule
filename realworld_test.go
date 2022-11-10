package rrule

import (
	"testing"
	"time"
)

func Test_real_event_for_fridays(t *testing.T) {
	// recurring weekly on friday with some recurring events removed.
	var value = "EXDATE;TZID=America/Los_Angeles:20180316T150000," +
		"20180323T150000,20180406T150000,20180427T150000," +
		"20180504T150000,20180511T150000,20180518T150000\n" +
		"RRULE:FREQ=WEEKLY;BYDAY=FR"

	rule, err := AssertToStringMatchesInput(value)
	if err != nil {
		t.Log("Failed to parse", err)
		t.Fail()
	}

	if rule.Frequency != WEEKLY {
		t.Log("Failed to parse frequency")
		t.Fail()
	}
	targetLa, _ := time.LoadLocation("America/Los_Angeles")

	expectedDates := []time.Time{
		time.Date(2018, 3, 16, 15, 0, 0, 0, targetLa),
		time.Date(2018, 3, 23, 15, 0, 0, 0, targetLa),
		time.Date(2018, 4, 6, 15, 0, 0, 0, targetLa),
		time.Date(2018, 4, 27, 15, 0, 0, 0, targetLa),
		time.Date(2018, 5, 4, 15, 0, 0, 0, targetLa),
		time.Date(2018, 5, 11, 15, 0, 0, 0, targetLa),
		time.Date(2018, 5, 18, 15, 0, 0, 0, targetLa),
	}

	for index, date := range expectedDates {
		if !rule.ExceptionsToRule[index].Equal(date) {
			t.Log("failed to match exception date", date)
			t.Fail()
		}
	}
}

func Test_FrequencyHourly(t *testing.T) {
	var DateStart = "DTSTART;TZID=America/New_York:20180902T090000\n"
	var RuleString = "RRULE:FREQ=HOURLY;COUNT=5"

	RuleShouldMatchDates(t,
		DateStart+RuleString,
		[]time.Time{
			time.Date(2018, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(2018, time.September, 2, 10, 0, 0, 0, targetLocation),
			time.Date(2018, time.September, 2, 11, 0, 0, 0, targetLocation),
			time.Date(2018, time.September, 2, 12, 0, 0, 0, targetLocation),
			time.Date(2018, time.September, 2, 13, 0, 0, 0, targetLocation),
		},
	)

}

func Test_FrequencyMinutely(t *testing.T) {
	var DateStart = "DTSTART;TZID=America/New_York:20180902T090000\n"
	var RuleString = "RRULE:FREQ=MINUTELY;COUNT=5"

	RuleShouldMatchDates(t,
		DateStart+RuleString,
		[]time.Time{
			time.Date(2018, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(2018, time.September, 2, 9, 1, 0, 0, targetLocation),
			time.Date(2018, time.September, 2, 9, 2, 0, 0, targetLocation),
			time.Date(2018, time.September, 2, 9, 3, 0, 0, targetLocation),
			time.Date(2018, time.September, 2, 9, 4, 0, 0, targetLocation),
		},
	)

}

func Test_FrequencySecondly(t *testing.T) {
	var DateStart = "DTSTART;TZID=America/New_York:20180902T090000\n"
	var RuleString = "RRULE:FREQ=SECONDLY;COUNT=5"

	RuleShouldMatchDates(t,
		DateStart+RuleString,
		[]time.Time{
			time.Date(2018, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(2018, time.September, 2, 9, 0, 1, 0, targetLocation),
			time.Date(2018, time.September, 2, 9, 0, 2, 0, targetLocation),
			time.Date(2018, time.September, 2, 9, 0, 3, 0, targetLocation),
			time.Date(2018, time.September, 2, 9, 0, 4, 0, targetLocation),
		},
	)

}

func Test_RealWorldFromGoogle_WithException(t *testing.T) {
	var DateStart = "DTSTART;TZID=America/New_York:20180902T090000\n"
	var Exceptions = "EXDATE;TZID=America/New_York:20180916T090000\n"
	var RuleString = "RRULE:FREQ=WEEKLY;COUNT=3;INTERVAL=2"

	RuleShouldMatchDates(t,
		DateStart+Exceptions+RuleString,
		[]time.Time{
			time.Date(2018, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(2018, time.September, 30, 9, 0, 0, 0, targetLocation),
			time.Date(2018, time.October, 14, 9, 0, 0, 0, targetLocation),
		},
	)

	RuleShouldMatchDates(t,
		DateStart+RuleString,
		[]time.Time{
			time.Date(2018, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(2018, time.September, 16, 9, 0, 0, 0, targetLocation),
			time.Date(2018, time.September, 30, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_RealWorldFromGoogle_WithException_DiffTimezone(t *testing.T) {
	var DateStart = "DTSTART;TZID=America/New_York:20180902T090000\n"
	var Exceptions = "EXDATE;TZID=America/Los_Angeles:20180916T060000\n"
	var RuleString = "RRULE:FREQ=WEEKLY;COUNT=3;INTERVAL=2"

	RuleShouldMatchDates(t,
		DateStart+Exceptions+RuleString,
		[]time.Time{
			time.Date(2018, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(2018, time.September, 30, 9, 0, 0, 0, targetLocation),
			time.Date(2018, time.October, 14, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_RealWorldFromGoogle_WithMultipleExceptions(t *testing.T) {
	var DateStart = "DTSTART;TZID=America/New_York:20180902T090000\n"
	var Exceptions = "EXDATE;TZID=America/New_York:20180916T090000,20180930T090000\n"
	var RuleString = "RRULE:FREQ=WEEKLY;COUNT=3;INTERVAL=2"

	RuleShouldMatchDates(t,
		DateStart+Exceptions+RuleString,
		[]time.Time{
			time.Date(2018, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(2018, time.October, 14, 9, 0, 0, 0, targetLocation),
			time.Date(2018, time.October, 28, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_Iterator_Limit(t *testing.T) {
	var s = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=HOURLY;INTERVAL=3;UNTIL=19970902T170000"

	rule, err := AssertToStringMatchesInput(s)
	if err != nil {
		t.Log("Failed to parse", err)
		t.Fail()
	}

	var event time.Time
	var index = 0

	iter := rule.Iterator().Limit(2)

	for iter.Step(&event) {
		index += 1
	}

	if index != 2 {
		t.Log("Failed to limit.")
		t.Fail()
	}
}

func Test_Iterator_Before(t *testing.T) {
	var s = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=HOURLY;INTERVAL=3;UNTIL=19970902T170000"

	rule, err := AssertToStringMatchesInput(s)
	if err != nil {
		t.Log("Failed to parse", err)
		t.Fail()
	}

	var event time.Time
	var index = 0

	// Without limits, results are as follows.
	// []time.Time{
	//   time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
	//   time.Date(1997, time.September, 2, 12, 0, 0, 0, targetLocation),
	//   time.Date(1997, time.September, 2, 15, 0, 0, 0, targetLocation),
	// },

	iter := rule.Iterator().Before(time.Date(1997, time.September, 2, 15, 0, 0, 0, targetLocation))

	for iter.Step(&event) {
		index += 1
	}

	if index != 2 {
		t.Log("Failed to restrict before.", index)
		t.Fail()
	}
}

func Test_Iterator_After(t *testing.T) {
	var s = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=HOURLY;INTERVAL=3;UNTIL=19970902T170000"

	rule, err := AssertToStringMatchesInput(s)
	if err != nil {
		t.Log("Failed to parse", err)
		t.Fail()
	}

	var event time.Time
	var index = 0

	// Without limits, results are as follows.
	// []time.Time{
	//   time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
	//   time.Date(1997, time.September, 2, 12, 0, 0, 0, targetLocation),
	//   time.Date(1997, time.September, 2, 15, 0, 0, 0, targetLocation),
	// },

	iter := rule.Iterator().After(time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation))

	for iter.Step(&event) {
		index += 1
	}

	if index != 2 {
		t.Log("Failed to restrict after.", index)
		t.Fail()
	}
}

func Test_ParseJustRule(t *testing.T) {
	var value string = "RRULE:FREQ=HOURLY;INTERVAL=3;COUNT=3"
	var dtStart = time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation)

	rule, err := AssertToStringMatchesInput(value)
	if err != nil {
		t.Log("Failed to parse", err)
		t.Fail()
	}

	rule.DtStart = dtStart

	var targetTimes []time.Time = []time.Time{
		time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 2, 12, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 2, 15, 0, 0, 0, targetLocation),
	}

	var index int = 0
	var event time.Time

	iter := rule.Iterator()

	for iter.Step(&event) {
		if !event.Equal(targetTimes[index]) {
			t.Log("Failed to match", event, targetTimes[index])
			t.Fail()
		} else {
			t.Log("Matched", event, targetTimes[index])
		}

		index += 1
	}
}

func Test_MonthlyFreqWithoutByMonthDay(t *testing.T) {
	var rule = "DTSTART;TZID=America/New_York:20200102T090000\nRRULE:FREQ=MONTHLY;COUNT=5"
	var expected = time.Date(2020, time.May, 02, 9, 0, 0, 0, targetLocation)

	r, err := Parse(rule)
	if err != nil {
		t.Fatal(err)
	}

	iter := r.Iterator()

	var out time.Time
	for iter.Step(&out) {
	}

	if !out.Equal(expected) {
		t.Fatal(out, expected)
	}
}

func Test_YearlyFreqWithoutByYearDay(t *testing.T) {
	var rule = "DTSTART;TZID=America/New_York:20200102T090000\nRRULE:FREQ=YEARLY;COUNT=5"
	var expected = time.Date(2024, time.January, 02, 9, 0, 0, 0, targetLocation)

	r, err := Parse(rule)
	if err != nil {
		t.Fatal(err)
	}

	iter := r.Iterator()

	var out time.Time
	for iter.Step(&out) {
	}

	if !out.Equal(expected) {
		t.Fatal(out, expected)
	}
}

func Test_IteratorHardLimit(t *testing.T) {
	var value = "DTSTART;TZID=America/New_York:20200102T090000\nRRULE:FREQ=DAILY"

	r, err := Parse(value)
	if err != nil {
		t.Fatal(err)
	}

	iter := r.Iterator().HardLimit(10)

	var tm time.Time
	for iter.Step(&tm) {
	}

	if !iter.IsHardLimitReached() {
		t.Fatal("Hard limit is not reached")
	}
}

func Benchmark_BasicParse(b *testing.B) {
	var value string = "RRULE:FREQ=HOURLY;INTERVAL=3;COUNT=3"

	for i := 0; i < b.N; i++ {
		_, err := Parse(value)
		if err != nil {
			panic(err)
		}

	}
}

func Benchmark_BasicParseAndIterate(b *testing.B) {
	var value string = "RRULE:FREQ=HOURLY;INTERVAL=3;COUNT=3"

	for i := 0; i < b.N; i++ {
		rule, err := Parse(value)
		if err != nil {
			panic(err)
		}

		iter := rule.Iterator()

		var event time.Time
		for iter.Step(&event) {

		}
	}
}

func Benchmark_ParseAndIterateMore(b *testing.B) {
	var value string = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=HOURLY;BYMINUTE=0,20,40;BYSECOND=10,20,30;UNTIL=19970904T090000"

	for i := 0; i < b.N; i++ {
		rule, err := Parse(value)
		if err != nil {
			panic(err)
		}

		iter := rule.Iterator()

		var event time.Time
		for iter.Step(&event) {

		}
	}
}

func Benchmark_ParseAndIterateTuesday(b *testing.B) {
	var value string = "DTSTART;TZID=America/New_York:19970310T090000\n" +
		"RRULE:FREQ=YEARLY;INTERVAL=2;COUNT=10;BYMONTH=1,2,3"

	for i := 0; i < b.N; i++ {
		rule, err := Parse(value)
		if err != nil {
			panic(err)
		}

		iter := rule.Iterator()

		var event time.Time
		for iter.Step(&event) {

		}
	}
}
