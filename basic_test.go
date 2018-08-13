package rrule

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

var targetLocation, _ = time.LoadLocation("America/New_York")

func AssertToStringMatchesInput(s string) (*RecurringRule, error) {
	rule, err := Parse(s)

	if err != nil {
		return rule, err
	}

	rule_checker, err := Parse(rule.String())
	if err != nil {
		return rule, err
	}

	if !rule.Equal(rule_checker) {
		return rule, errors.New(
			fmt.Sprintf(
				"Rendered String doesn't match input: \n%s\n%s",
				s,
				rule.String(),
			),
		)
	}

	return rule, nil
}

func RuleShouldMatchDates(t *testing.T, s string, targetDates []time.Time) {
	rule, err := AssertToStringMatchesInput(s)
	if err != nil {
		t.Log("Failed to parse", err)
		t.Fail()
	}

	var event time.Time
	var index = 0

	iter := rule.Iterator()

	for iter.Step(&event) {
		if index >= len(targetDates) {
			return
		}
		if event.Equal(targetDates[index]) == false {
			t.Log(
				fmt.Sprintf(
					"Failed to match %s and %s",
					event,
					targetDates[index],
				),
			)
			t.Fail()
		} else {
			t.Log(
				fmt.Sprintf(
					"Matched %s and %s",
					event,
					targetDates[index],
				),
			)
		}
		index += 1
	}

	if iter.ReturnCounter != len(targetDates) {
		t.Log("Failed to return enough results.", iter.ReturnCounter, len(targetDates))
		t.Fail()
	}

	t.Log("----")
}

func Test_DailyFor10(t *testing.T) {
	// Daily for 10 occurrences:

	//  DTSTART;TZID=America/New_York:19970902T090000
	//  RRULE:FREQ=DAILY;COUNT=10

	//  ==> (1997 9:00 AM EDT) September 2-11

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970902T090000\n"+
			"RRULE:FREQ=DAILY;COUNT=10",
		[]time.Time{
			time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 3, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 4, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 5, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 6, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 7, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 8, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 9, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 11, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_DailyUntil(t *testing.T) {
	// Daily until December 24, 1997:

	//  DTSTART;TZID=America/New_York:19970902T090000
	//  RRULE:FREQ=DAILY;UNTIL=19971224T000000Z

	//  ==> (1997 9:00 AM EDT) September 2-30;October 1-25
	//      (1997 9:00 AM EST) October 26-31;November 1-30;December 1-23

	var times = []time.Time{}

	var temp time.Time = time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation)

	for i := 0; i < (28 + 31 + 30 + 23 + 1); i += 1 {
		times = append(times, temp)
		temp = temp.AddDate(0, 0, 1)
	}

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970902T090000\n"+
			"RRULE:FREQ=DAILY;UNTIL=19971224T000000Z",
		times,
	)

	if !times[len(times)-1].Equal(time.Date(1997, time.December, 23, 9, 0, 0, 0, targetLocation)) {
		t.Log("Incorrect times generated.")
		t.Fail()
	}
}

func Test_Daily_EveryOtherDayForever(t *testing.T) {
	// Every other day - forever:

	//  DTSTART;TZID=America/New_York:19970902T090000
	//  RRULE:FREQ=DAILY;INTERVAL=2

	//  ==> (1997 9:00 AM EDT) September 2,4,6,8...24,26,28,30;
	//                         October 2,4,6...20,22,24
	//      (1997 9:00 AM EST) October 26,28,30;
	//                         November 1,3,5,7...25,27,29;
	// 	December 1,3,...

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970902T090000\n"+
			"RRULE:FREQ=DAILY;INTERVAL=2",
		[]time.Time{
			time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 4, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 6, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 8, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 12, 9, 0, 0, 0, targetLocation),
		},
	)

}

// https://tools.ietf.org/html/rfc5545#page-124
func Test_Interval_And_Count(t *testing.T) {
	//       Every 10 days, 5 occurrences:
	//
	//         DTSTART;TZID=America/New_York:19970902T090000
	//         RRULE:FREQ=DAILY;INTERVAL=10;COUNT=5
	//
	//       ==> (1997 9:00 AM EDT) September 2,12,22;
	//                              October 2,12

	var value = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=DAILY;INTERVAL=10;COUNT=5"

	var targetDates []time.Time = []time.Time{
		time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 12, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 22, 9, 0, 0, 0, targetLocation),

		time.Date(1997, time.October, 2, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.October, 12, 9, 0, 0, 0, targetLocation),
	}

	rule, err := AssertToStringMatchesInput(value)
	if err != nil {
		t.Log("Failed to parse", err)
		t.Fail()
	}

	var event time.Time
	var index = 0

	iter := rule.Iterator()

	for iter.Step(&event) {
		if event.Equal(targetDates[index]) == false {
			t.Log(
				fmt.Sprintf(
					"Failed to match %s and %s",
					event,
					targetDates[index],
				),
			)
			t.Fail()
		}
		index += 1
	}
}

// https://tools.ietf.org/html/rfc5545#page-124
func Test_EveryDayInJanuary_ForThreeYears(t *testing.T) {
	// Every day in January, for 3 years:
	//
	//  DTSTART;TZID=America/New_York:19980101T090000
	//
	//  RRULE:FREQ=YEARLY;UNTIL=20000131T140000Z;
	//   BYMONTH=1;BYDAY=SU,MO,TU,WE,TH,FR,SA
	//  or
	//  RRULE:FREQ=DAILY;UNTIL=20000131T140000Z;BYMONTH=1
	//
	//  ==> (1998 9:00 AM EST)January 1-31
	//      (1999 9:00 AM EST)January 1-31
	//      (2000 9:00 AM EST)January 1-31

	var values = []string{
		"DTSTART;TZID=America/New_York:19980101T090000\n" +
			"RRULE:FREQ=YEARLY;UNTIL=20000131T140000Z;" +
			"BYMONTH=1;BYDAY=SU,MO,TU,WE,TH,FR,SA",
		"DTSTART;TZID=America/New_York:19980101T090000\n" +
			"RRULE:FREQ=DAILY;UNTIL=20000131T140000Z;BYMONTH=1",
	}

	var targetDates = []time.Time{}

	for _, year := range []int{1998, 1999, 2000} {
		for day := 1; day < 32; day += 1 {
			targetDates = append(
				targetDates,
				time.Date(
					year,
					time.Month(1),
					day,
					9,
					0,
					0,
					0,
					targetLocation,
				),
			)
		}
	}

	for _, v := range values {
		rule, err := AssertToStringMatchesInput(v)

		if err != nil {
			t.Log(
				fmt.Sprintf(
					"Failed to parse rule: %s",
					v,
				),
			)
			t.Fail()
		}

		iter := rule.Iterator()

		var event time.Time
		var index int = 0

		for iter.Step(&event) {
			if event.Equal(targetDates[index]) == false {
				t.Log("Failed to match", event, targetDates[index])
				t.Fail()
			}
			index += 1
		}

		if iter.ReturnCounter != len(targetDates) {
			t.Log("Not enough results", iter.ReturnCounter, len(targetDates))
			t.Fail()
		}
	}
}

// https://tools.ietf.org/html/rfc5545#page-124
func Test_Weekly_for_10_Count(t *testing.T) {
	// Weekly for 10 occurrences:
	//
	//   DTSTART;TZID=America/New_York:19970902T090000
	//   RRULE:FREQ=WEEKLY;COUNT=10
	//
	//   ==> (1997 9:00 AM EDT) September 2,9,16,23,30;October 7,14,21
	//       (1997 9:00 AM EST) October 28;November 4

	var value = "DTSTART;TZID=America/New_York:19970902T090000\nRRULE:FREQ=WEEKLY;COUNT=10"

	rule, err := AssertToStringMatchesInput(value)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	var targetDates = []time.Time{
		time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 9, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 16, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 23, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 30, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.October, 7, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.October, 14, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.October, 21, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.October, 28, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.November, 4, 9, 0, 0, 0, targetLocation),
	}

	var event time.Time
	var index = 0

	iter := rule.Iterator()

	for iter.Step(&event) {
		if event.Equal(targetDates[index]) != true {
			t.Fail()
		}
		index += 1
	}
}

func Test_TU_TH_for_Five_Weeks(t *testing.T) {
	// Weekly on Tuesday and Thursday for five weeks:
	//
	//  DTSTART;TZID=America/New_York:19970902T090000
	//  RRULE:FREQ=WEEKLY;UNTIL=19971007T000000Z;WKST=SU;BYDAY=TU,TH
	//
	//  or
	//
	//  RRULE:FREQ=WEEKLY;COUNT=10;WKST=SU;BYDAY=TU,TH
	//
	//  ==> (1997 9:00 AM EDT) September 2,4,9,11,16,18,23,25,30;
	//                         October 2

	var values = []string{
		"DTSTART;TZID=America/New_York:19970902T090000\n" +
			"RRULE:FREQ=WEEKLY;UNTIL=19971007T000000Z;WKST=SU;BYDAY=TU,TH",
		"DTSTART;TZID=America/New_York:19970902T090000\n" +
			"RRULE:FREQ=WEEKLY;COUNT=10;WKST=SU;BYDAY=TU,TH",
	}

	var targetDates = []time.Time{
		time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 4, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 9, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 11, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 16, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 18, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 23, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 25, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 30, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.October, 2, 9, 0, 0, 0, targetLocation),
	}

	for _, v := range values {
		rule, err := AssertToStringMatchesInput(v)

		if err != nil {
			t.Log(
				fmt.Sprintf(
					"Failed to parse rule: %s",
					v,
				),
			)
			t.Fail()
		}

		var event time.Time
		var index = 0

		iter := rule.Iterator()

		for iter.Step(&event) {
			if event.Equal(targetDates[index]) != true {
				t.Fail()
			}
			index += 1
		}

		if index != len(targetDates) {
			t.Log("Incorrect iter len")
			t.Fail()
		}
	}
}

func Test_every_other_complex(t *testing.T) {
	// Every other week on Monday, Wednesday, and Friday until December
	// 24, 1997, starting on Monday, September 1, 1997:
	//
	//  DTSTART;TZID=America/New_York:19970901T090000
	//  RRULE:FREQ=WEEKLY;INTERVAL=2;UNTIL=19971224T000000Z;WKST=SU;
	//   BYDAY=MO,WE,FR
	//
	//  ==> (1997 9:00 AM EDT) September 1,3,5,15,17,19,29;
	//                         October 1,3,13,15,17
	//      (1997 9:00 AM EST) October 27,29,31;
	//                         November 10,12,14,24,26,28;
	//                         December 8,10,12,22

	RuleShouldMatchDates(
		t,
		"DTSTART;TZID=America/New_York:19970901T090000\n"+
			"RRULE:FREQ=WEEKLY;INTERVAL=2;UNTIL=19971224T000000Z;WKST=SU;"+
			"BYDAY=MO,WE,FR",
		[]time.Time{
			time.Date(1997, time.September, 1, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 3, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 5, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 15, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 17, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 19, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 29, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 1, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 3, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 13, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 15, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 17, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 27, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 29, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 31, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 12, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 14, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 24, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 26, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 28, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 8, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 12, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 22, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_Monthly_First_Friday(t *testing.T) {
	//  Monthly on the first Friday until December 24, 1997:

	//  DTSTART;TZID=America/New_York:19970905T090000
	//  RRULE:FREQ=MONTHLY;UNTIL=19971224T000000Z;BYDAY=1FR

	//  ==> (1997 9:00 AM EDT) September 5; October 3
	//      (1997 9:00 AM EST) November 7; December 5

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970905T090000\n"+
			"RRULE:FREQ=MONTHLY;UNTIL=19971224T000000Z;BYDAY=1FR",
		[]time.Time{
			time.Date(1997, time.September, 5, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 3, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 7, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 5, 9, 0, 0, 0, targetLocation),
		},
	)

}

func Test_Monthly_on_the_third_to_last_day_of_month(t *testing.T) {
	// Monthly on the third-to-the-last day of the month, forever:

	//  DTSTART;TZID=America/New_York:19970928T090000
	//  RRULE:FREQ=MONTHLY;BYMONTHDAY=-3

	//  ==> (1997 9:00 AM EDT) September 28
	//      (1997 9:00 AM EST) October 29;November 28;December 29
	//      (1998 9:00 AM EST) January 29;February 26
	//      ...

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970928T090000\n"+
			"RRULE:FREQ=MONTHLY;BYMONTHDAY=-3",
		[]time.Time{
			time.Date(1997, time.September, 28, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 29, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 28, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 29, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 29, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.February, 26, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_BySetPos(t *testing.T) {
	// The third instance into the month of one of Tuesday, Wednesday, or
	// Thursday, for the next 3 months:

	//  DTSTART;TZID=America/New_York:19970904T090000
	//  RRULE:FREQ=MONTHLY;COUNT=3;BYDAY=TU,WE,TH;BYSETPOS=3

	//  ==> (1997 9:00 AM EDT) September 4;October 7
	//      (1997 9:00 AM EST) November 6

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970904T090000\n"+
			"RRULE:FREQ=MONTHLY;COUNT=3;BYDAY=TU,WE,TH;BYSETPOS=3",
		[]time.Time{
			time.Date(1997, time.September, 4, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 7, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 6, 9, 0, 0, 0, targetLocation),
		},
	)

	// The second-to-last weekday of the month:

	//  DTSTART;TZID=America/New_York:19970929T090000
	//  RRULE:FREQ=MONTHLY;BYDAY=MO,TU,WE,TH,FR;BYSETPOS=-2

	//  ==> (1997 9:00 AM EDT) September 29
	//      (1997 9:00 AM EST) October 30;November 27;December 30
	//      (1998 9:00 AM EST) January 29;February 26;March 30
	//      ...

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970929T090000\n"+
			"RRULE:FREQ=MONTHLY;BYDAY=MO,TU,WE,TH,FR;BYSETPOS=-2",
		[]time.Time{
			time.Date(1997, time.September, 29, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 30, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 27, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 30, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 29, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.February, 26, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 30, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_ByYearDay(t *testing.T) {
	// Every third year on the 1st, 100th, and 200th day for 10
	// occurrences:

	//  DTSTART;TZID=America/New_York:19970101T090000
	//  RRULE:FREQ=YEARLY;INTERVAL=3;COUNT=10;BYYEARDAY=1,100,200

	//  ==> (1997 9:00 AM EST) January 1
	//      (1997 9:00 AM EDT) April 10;July 19
	//      (2000 9:00 AM EST) January 1
	//      (2000 9:00 AM EDT) April 9;July 18
	//      (2003 9:00 AM EST) January 1
	//      (2003 9:00 AM EDT) April 10;July 19
	//      (2006 9:00 AM EST) January 1

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970101T090000\n"+
			"RRULE:FREQ=YEARLY;INTERVAL=3;COUNT=10;BYYEARDAY=1,100,200",
		[]time.Time{
			time.Date(1997, time.January, 1, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.April, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.July, 19, 9, 0, 0, 0, targetLocation),
			time.Date(2000, time.January, 1, 9, 0, 0, 0, targetLocation),
			time.Date(2000, time.April, 9, 9, 0, 0, 0, targetLocation),
			time.Date(2000, time.July, 18, 9, 0, 0, 0, targetLocation),
			time.Date(2003, time.January, 1, 9, 0, 0, 0, targetLocation),
			time.Date(2003, time.April, 10, 9, 0, 0, 0, targetLocation),
			time.Date(2003, time.July, 19, 9, 0, 0, 0, targetLocation),
			time.Date(2006, time.January, 1, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_FirstandLast_SundayofMonth(t *testing.T) {
	// Every other month on the first and last Sunday of the month for 10
	// occurrences:

	//  DTSTART;TZID=America/New_York:19970907T090000
	//  RRULE:FREQ=MONTHLY;INTERVAL=2;COUNT=10;BYDAY=1SU,-1SU

	//  ==> (1997 9:00 AM EDT) September 7,28
	//      (1997 9:00 AM EST) November 2,30
	//      (1998 9:00 AM EST) January 4,25;March 1,29
	//      (1998 9:00 AM EDT) May 3,31

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970907T090000\n"+
			"RRULE:FREQ=MONTHLY;INTERVAL=2;COUNT=10;BYDAY=1SU,-1SU",
		[]time.Time{
			time.Date(1997, time.September, 7, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 28, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 30, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 4, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 25, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 1, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 29, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.May, 3, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.May, 31, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_Monthly_secondtolast_monday(t *testing.T) {
	// Monthly on the second-to-last Monday of the month for 6 months:

	//  DTSTART;TZID=America/New_York:19970922T090000
	//  RRULE:FREQ=MONTHLY;COUNT=6;BYDAY=-2MO

	//  ==> (1997 9:00 AM EDT) September 22;October 20
	//      (1997 9:00 AM EST) November 17;December 22
	//      (1998 9:00 AM EST) January 19;February 16

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970922T090000\n"+
			"RRULE:FREQ=MONTHLY;COUNT=6;BYDAY=-2MO",
		[]time.Time{
			time.Date(1997, time.September, 22, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 20, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 17, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 22, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 19, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.February, 16, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_FirstAndLastDayOfMonth(t *testing.T) {
	// Monthly on the first and last day of the month for 10 occurrences:

	//  DTSTART;TZID=America/New_York:19970930T090000
	//  RRULE:FREQ=MONTHLY;COUNT=10;BYMONTHDAY=1,-1

	//  ==> (1997 9:00 AM EDT) September 30;October 1
	//      (1997 9:00 AM EST) October 31;November 1,30;December 1,31
	//      (1998 9:00 AM EST) January 1,31;February 1

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970930T090000\n"+
			"RRULE:FREQ=MONTHLY;COUNT=10;BYMONTHDAY=1,-1",
		[]time.Time{
			time.Date(1997, time.September, 30, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 1, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 31, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 1, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 30, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 1, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 31, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 1, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 31, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.February, 1, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_LongInterval(t *testing.T) {
	// Every 18 months on the 10th thru 15th of the month for 10
	// occurrences:

	//  DTSTART;TZID=America/New_York:19970910T090000
	//  RRULE:FREQ=MONTHLY;INTERVAL=18;COUNT=10;BYMONTHDAY=10,11,12,
	//   13,14,15

	//  ==> (1997 9:00 AM EDT) September 10,11,12,13,14,15
	//      (1999 9:00 AM EST) March 10,11,12,13

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970910T090000\n"+
			"RRULE:FREQ=MONTHLY;INTERVAL=18;COUNT=10;BYMONTHDAY=10,11,12,13,14,15",
		[]time.Time{
			time.Date(1997, time.September, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 11, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 12, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 13, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 14, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 15, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.March, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.March, 11, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.March, 12, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.March, 13, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_YearlyWithoutSpecifics(t *testing.T) {
	// Yearly in June and July for 10 occurrences:

	//  DTSTART;TZID=America/New_York:19970610T090000
	//  RRULE:FREQ=YEARLY;COUNT=10;BYMONTH=6,7

	//  ==> (1997 9:00 AM EDT) June 10;July 10
	//      (1998 9:00 AM EDT) June 10;July 10
	//      (1999 9:00 AM EDT) June 10;July 10
	//      (2000 9:00 AM EDT) June 10;July 10
	//      (2001 9:00 AM EDT) June 10;July 10

	//    Note: Since none of the BYDAY, BYMONTHDAY, or BYYEARDAY
	//    components are specified, the day is gotten from "DTSTART".

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970610T090000\n"+
			"RRULE:FREQ=YEARLY;COUNT=10;BYMONTH=6,7",
		[]time.Time{
			time.Date(1997, time.June, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.July, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.June, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.July, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.June, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.July, 10, 9, 0, 0, 0, targetLocation),
			time.Date(2000, time.June, 10, 9, 0, 0, 0, targetLocation),
			time.Date(2000, time.July, 10, 9, 0, 0, 0, targetLocation),
			time.Date(2001, time.June, 10, 9, 0, 0, 0, targetLocation),
			time.Date(2001, time.July, 10, 9, 0, 0, 0, targetLocation),
		},
	)

}

func Test_Every20thMonday(t *testing.T) {
	// Every 20th Monday of the year, forever:

	//  DTSTART;TZID=America/New_York:19970519T090000
	//  RRULE:FREQ=YEARLY;BYDAY=20MO

	//  ==> (1997 9:00 AM EDT) May 19
	//      (1998 9:00 AM EDT) May 18
	//      (1999 9:00 AM EDT) May 17
	//      ...

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970519T090000\n"+
			"RRULE:FREQ=YEARLY;BYDAY=20MO",
		[]time.Time{
			time.Date(1997, time.May, 19, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.May, 18, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.May, 17, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_Monthly_2nd_15th_for_10Count(t *testing.T) {
	// Monthly on the 2nd and 15th of the month for 10 occurrences:

	//  DTSTART;TZID=America/New_York:19970902T090000
	//  RRULE:FREQ=MONTHLY;COUNT=10;BYMONTHDAY=2,15

	//  ==> (1997 9:00 AM EDT) September 2,15;October 2,15
	//      (1997 9:00 AM EST) November 2,15;December 2,15
	//      (1998 9:00 AM EST) January 2,15

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970902T090000\n"+
			"RRULE:FREQ=MONTHLY;COUNT=10;BYMONTHDAY=2,15",
		[]time.Time{
			time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 15, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 15, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 15, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 15, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 15, 9, 0, 0, 0, targetLocation),
		},
	)

}

func Test_EveryTuesday_EveryOtherMonth(t *testing.T) {
	// Every Tuesday, every other month:

	//  DTSTART;TZID=America/New_York:19970902T090000
	//  RRULE:FREQ=MONTHLY;INTERVAL=2;BYDAY=TU

	//  ==> (1997 9:00 AM EDT) September 2,9,16,23,30
	//      (1997 9:00 AM EST) November 4,11,18,25
	//      (1998 9:00 AM EST) January 6,13,20,27;March 3,10,17,24,31
	//      ...

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970902T090000\n"+
			"RRULE:FREQ=MONTHLY;INTERVAL=2;BYDAY=TU",
		[]time.Time{
			time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 9, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 16, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 23, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 30, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 4, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 11, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 18, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 25, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 6, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 13, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 20, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 27, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 3, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 17, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 24, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 31, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_EveryOtherYear_onJFM(t *testing.T) {
	// Every other year on January, February, and March for 10
	// occurrences:

	//  DTSTART;TZID=America/New_York:19970310T090000
	//  RRULE:FREQ=YEARLY;INTERVAL=2;COUNT=10;BYMONTH=1,2,3

	//  ==> (1997 9:00 AM EST) March 10
	//      (1999 9:00 AM EST) January 10;February 10;March 10
	//      (2001 9:00 AM EST) January 10;February 10;March 10
	//      (2003 9:00 AM EST) January 10;February 10;March 10

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970310T090000\n"+
			"RRULE:FREQ=YEARLY;INTERVAL=2;COUNT=10;BYMONTH=1,2,3",
		[]time.Time{
			time.Date(1997, time.March, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.January, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.February, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.March, 10, 9, 0, 0, 0, targetLocation),
			time.Date(2001, time.January, 10, 9, 0, 0, 0, targetLocation),
			time.Date(2001, time.February, 10, 9, 0, 0, 0, targetLocation),
			time.Date(2001, time.March, 10, 9, 0, 0, 0, targetLocation),
			time.Date(2003, time.January, 10, 9, 0, 0, 0, targetLocation),
			time.Date(2003, time.February, 10, 9, 0, 0, 0, targetLocation),
			time.Date(2003, time.March, 10, 9, 0, 0, 0, targetLocation),
		},
	)

}

func Test_EveryThursdayInMarch(t *testing.T) {
	// Every Thursday in March, forever:

	//  DTSTART;TZID=America/New_York:19970313T090000
	//  RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=TH

	//  ==> (1997 9:00 AM EST) March 13,20,27
	//      (1998 9:00 AM EST) March 5,12,19,26
	//      (1999 9:00 AM EST) March 4,11,18,25
	//      ...

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970313T090000\n"+
			"RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=TH",
		[]time.Time{
			time.Date(1997, time.March, 13, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.March, 20, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.March, 27, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 5, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 12, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 19, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 26, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.March, 4, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.March, 11, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.March, 18, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.March, 25, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_EveryT_but_OnlyJuneJulyAugust(t *testing.T) {
	// Every Thursday, but only during June, July, and August, forever:

	//  DTSTART;TZID=America/New_York:19970605T090000
	//  RRULE:FREQ=YEARLY;BYDAY=TH;BYMONTH=6,7,8

	//  ==> (1997 9:00 AM EDT) June 5,12,19,26;July 3,10,17,24,31;
	//                         August 7,14,21,28
	//      (1998 9:00 AM EDT) June 4,11,18,25;July 2,9,16,23,30;
	//                         August 6,13,20,27
	//      (1999 9:00 AM EDT) June 3,10,17,24;July 1,8,15,22,29;
	//                         August 5,12,19,26
	//      ...

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970605T090000\n"+
			"RRULE:FREQ=YEARLY;BYDAY=TH;BYMONTH=6,7,8",
		[]time.Time{
			time.Date(1997, time.June, 5, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.June, 12, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.June, 19, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.June, 26, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.July, 3, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.July, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.July, 17, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.July, 24, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.July, 31, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.August, 7, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.August, 14, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.August, 21, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.August, 28, 9, 0, 0, 0, targetLocation),

			time.Date(1998, time.June, 4, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.June, 11, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.June, 18, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.June, 25, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.July, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.July, 9, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.July, 16, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.July, 23, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.July, 30, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.August, 6, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.August, 13, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.August, 20, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.August, 27, 9, 0, 0, 0, targetLocation),

			time.Date(1999, time.June, 3, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.June, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.June, 17, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.June, 24, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.July, 1, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.July, 8, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.July, 15, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.July, 22, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.July, 29, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.August, 5, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.August, 12, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.August, 19, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.August, 26, 9, 0, 0, 0, targetLocation),
		},
	)

}

func Test_FridayThe13th(t *testing.T) {
	// Every Friday the 13th, forever:

	//  DTSTART;TZID=America/New_York:19970902T090000
	//  EXDATE;TZID=America/New_York:19970902T090000
	//  RRULE:FREQ=MONTHLY;BYDAY=FR;BYMONTHDAY=13

	//  ==> (1998 9:00 AM EST) February 13;March 13;November 13
	//      (1999 9:00 AM EDT) August 13
	//      (2000 9:00 AM EDT) October 13
	//      ...

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970902T090000\n"+
			"EXDATE;TZID=America/New_York:19970902T090000\n"+
			"RRULE:FREQ=MONTHLY;BYDAY=FR;BYMONTHDAY=13",
		[]time.Time{
			time.Date(1998, time.February, 13, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 13, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.November, 13, 9, 0, 0, 0, targetLocation),
			time.Date(1999, time.August, 13, 9, 0, 0, 0, targetLocation),
			time.Date(2000, time.October, 13, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_FirstSaturdayAfterFirstSunday(t *testing.T) {
	// The first Saturday that follows the first Sunday of the month,
	// forever:

	//  DTSTART;TZID=America/New_York:19970913T090000
	//  RRULE:FREQ=MONTHLY;BYDAY=SA;BYMONTHDAY=7,8,9,10,11,12,13

	//  ==> (1997 9:00 AM EDT) September 13;October 11
	//      (1997 9:00 AM EST) November 8;December 13
	//      (1998 9:00 AM EST) January 10;February 7;March 7
	//      (1998 9:00 AM EDT) April 11;May 9;June 13...
	//      ...

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970913T090000\n"+
			"RRULE:FREQ=MONTHLY;BYDAY=SA;BYMONTHDAY=7,8,9,10,11,12,13",
		[]time.Time{
			time.Date(1997, time.September, 13, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.October, 11, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.November, 8, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.December, 13, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.January, 10, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.February, 7, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.March, 7, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.April, 11, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.May, 9, 9, 0, 0, 0, targetLocation),
			time.Date(1998, time.June, 13, 9, 0, 0, 0, targetLocation),
		},
	)

}

func Test_Election_Day_US(t *testing.T) {
	// Every 4 years, the first Tuesday after a Monday in November,
	// forever (U.S. Presidential Election day):

	//  DTSTART;TZID=America/New_York:19961105T090000
	//  RRULE:FREQ=YEARLY;INTERVAL=4;BYMONTH=11;BYDAY=TU;
	//   BYMONTHDAY=2,3,4,5,6,7,8

	//   ==> (1996 9:00 AM EST) November 5
	//       (2000 9:00 AM EST) November 7
	//       (2004 9:00 AM EST) November 2
	//       ...

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19961105T090000\n"+
			"RRULE:FREQ=YEARLY;INTERVAL=4;BYMONTH=11;BYDAY=TU;"+
			"BYMONTHDAY=2,3,4,5,6,7,8",
		[]time.Time{
			time.Date(1996, time.November, 5, 9, 0, 0, 0, targetLocation),
			time.Date(2000, time.November, 7, 9, 0, 0, 0, targetLocation),
			time.Date(2004, time.November, 2, 9, 0, 0, 0, targetLocation),
		},
	)
}

func Test_WeekNo(t *testing.T) {
	// Monday of week number 20 (where the default start of the week is
	// Monday), forever:

	//  DTSTART;TZID=America/New_York:19970512T090000
	//  RRULE:FREQ=YEARLY;BYWEEKNO=20;BYDAY=MO

	//  ==> (1997 9:00 AM EDT) May 12
	//      (1998 9:00 AM EDT) May 11
	//      (1999 9:00 AM EDT) May 17
	//      ...

	var value = "DTSTART;TZID=America/New_York:19970512T090000\n" +
		"RRULE:FREQ=YEARLY;BYWEEKNO=20;BYDAY=MO"

	RuleShouldMatchDates(t, value, []time.Time{
		time.Date(1997, time.May, 12, 9, 0, 0, 0, targetLocation),
		time.Date(1998, time.May, 11, 9, 0, 0, 0, targetLocation),
		time.Date(1999, time.May, 17, 9, 0, 0, 0, targetLocation),
	})
}

func Test_Dont_Return_Invalid_Date(t *testing.T) {
	// An example where an invalid date (i.e., February 30) is ignored.

	//  DTSTART;TZID=America/New_York:20070115T090000
	//  RRULE:FREQ=MONTHLY;BYMONTHDAY=15,30;COUNT=5

	//  ==> (2007 EST) January 15,30
	//      (2007 EST) February 15
	//      (2007 EDT) March 15,30

	var value = "DTSTART;TZID=America/New_York:20070115T090000\n" +
		"RRULE:FREQ=MONTHLY;BYMONTHDAY=15,30;COUNT=5"

	RuleShouldMatchDates(t, value, []time.Time{
		time.Date(2007, time.January, 15, 9, 0, 0, 0, targetLocation),
		time.Date(2007, time.January, 30, 9, 0, 0, 0, targetLocation),
		time.Date(2007, time.February, 15, 9, 0, 0, 0, targetLocation),
		time.Date(2007, time.March, 15, 9, 0, 0, 0, targetLocation),
		time.Date(2007, time.March, 30, 9, 0, 0, 0, targetLocation),
	})

}

func Test_WKST_Change(t *testing.T) {
	// An example where the days generated makes a difference because of
	// WKST:
	//
	//  DTSTART;TZID=America/New_York:19970805T090000
	//  RRULE:FREQ=WEEKLY;INTERVAL=2;COUNT=4;BYDAY=TU,SU;WKST=MO
	//
	//  ==> (1997 EDT) August 5,10,19,24
	//
	// changing only WKST from MO to SU, yields different results...
	//
	//  DTSTART;TZID=America/New_York:19970805T090000
	//  RRULE:FREQ=WEEKLY;INTERVAL=2;COUNT=4;BYDAY=TU,SU;WKST=SU
	//
	//  ==> (1997 EDT) August 5,17,19,31

	var value = "DTSTART;TZID=America/New_York:19970805T090000\n" +
		"RRULE:FREQ=WEEKLY;INTERVAL=2;COUNT=4;BYDAY=TU,SU;"

	var value_wkst_monday = value + "WKST=MO"
	var value_wkst_sunday = value + "WKST=SU"

	var targetTimesWkstMonday []time.Time = []time.Time{
		time.Date(1997, 8, 5, 9, 0, 0, 0, targetLocation),
		time.Date(1997, 8, 10, 9, 0, 0, 0, targetLocation),
		time.Date(1997, 8, 19, 9, 0, 0, 0, targetLocation),
		time.Date(1997, 8, 24, 9, 0, 0, 0, targetLocation),
	}

	var targetTimesWkstSunday []time.Time = []time.Time{
		time.Date(1997, 8, 5, 9, 0, 0, 0, targetLocation),
		time.Date(1997, 8, 17, 9, 0, 0, 0, targetLocation),
		time.Date(1997, 8, 19, 9, 0, 0, 0, targetLocation),
		time.Date(1997, 8, 31, 9, 0, 0, 0, targetLocation),
	}

	ruleMonday, err := AssertToStringMatchesInput(value_wkst_monday)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	var iter *RecurrenceIterator = ruleMonday.Iterator()
	var event time.Time
	var index = 0
	for iter.Step(&event) {
		if event.Equal(targetTimesWkstMonday[index]) == false {
			t.Log(
				fmt.Sprintf(
					"Failed match [%s] %s and [%s] %s",
					event.Weekday(),
					event,
					targetTimesWkstMonday[index].Weekday(),
					targetTimesWkstMonday[index],
				),
			)
			t.Fail()
		}
		index += 1
	}

	ruleSunday, err := AssertToStringMatchesInput(value_wkst_sunday)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	iter = ruleSunday.Iterator()
	index = 0
	for iter.Step(&event) {
		if event.Equal(targetTimesWkstSunday[index]) == false {
			t.Log(
				fmt.Sprintf(
					"Failed match [%s] %s and [%s] %s",
					event.Weekday(),
					event,
					targetTimesWkstSunday[index].Weekday(),
					targetTimesWkstSunday[index],
				),
			)
			t.Fail()
		}
		index += 1
	}

}

func Test_TimeOfDay_WithUntil(t *testing.T) {
	// Every 3 hours from 9:00 AM to 5:00 PM on a specific day:

	//  DTSTART;TZID=America/New_York:19970902T090000
	//  RRULE:FREQ=HOURLY;INTERVAL=3;UNTIL=19970902T170000Z

	// This is incorrect in the RFC, in the rfc errata I found
	// a correction to UNTIL=19970902T210000Z (ahead of the original)

	//  ==> (September 2, 1997 EDT) 09:00,12:00,15:00

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970902T090000\n"+
			"RRULE:FREQ=HOURLY;INTERVAL=3;UNTIL=19970902T210000Z",
		[]time.Time{
			time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, 12, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, 15, 0, 0, 0, targetLocation),
		},
	)
}

func Test_TimeOfDay_WithCount(t *testing.T) {
	// Every 15 minutes for 6 occurrences:

	//  DTSTART;TZID=America/New_York:19970902T090000
	//  RRULE:FREQ=MINUTELY;INTERVAL=15;COUNT=6

	//  ==> (September 2, 1997 EDT) 09:00,09:15,09:30,09:45,10:00,10:15

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970902T090000\n"+
			"RRULE:FREQ=MINUTELY;INTERVAL=15;COUNT=6",
		[]time.Time{
			time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, 9, 15, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, 9, 30, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, 9, 45, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 15, 0, 0, targetLocation),
		},
	)

	// Every hour and a half for 4 occurrences:

	//  DTSTART;TZID=America/New_York:19970902T090000
	//  RRULE:FREQ=MINUTELY;INTERVAL=90;COUNT=4

	//  ==> (September 2, 1997 EDT) 09:00,10:30;12:00;13:30

	RuleShouldMatchDates(t,
		"DTSTART;TZID=America/New_York:19970902T090000\n"+
			"RRULE:FREQ=MINUTELY;INTERVAL=90;COUNT=4",
		[]time.Time{
			time.Date(1997, time.September, 2, 9, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 30, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, 12, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, 13, 30, 0, 0, targetLocation),
		},
	)
}

func Test_ByMinute(t *testing.T) {
	// Every 20 minutes from 9:00 AM to 4:40 PM every day:

	//  DTSTART;TZID=America/New_York:19970902T090000
	//  RRULE:FREQ=DAILY;BYHOUR=9,10,11,12,13,14,15,16;BYMINUTE=0,20,40
	//  or
	//  RRULE:FREQ=MINUTELY;INTERVAL=20;BYHOUR=9,10,11,12,13,14,15,16

	//  ==> (September 2, 1997 EDT) 9:00,9:20,9:40,10:00,10:20,
	//                              ... 16:00,16:20,16:40
	//      (September 3, 1997 EDT) 9:00,9:20,9:40,10:00,10:20,
	//                              ...16:00,16:20,16:40
	//      ...

	var s string = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=DAILY;BYHOUR=9,10,11,12,13,14,15,16;BYMINUTE=0,20,40"

	var s2 string = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=MINUTELY;INTERVAL=20;BYHOUR=9,10,11,12,13,14,15,16"

	var times []time.Time

	for i := 9; i < 17; i += 1 {
		times = append(times, []time.Time{
			time.Date(1997, time.September, 2, i, 0, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, i, 20, 0, 0, targetLocation),
			time.Date(1997, time.September, 2, i, 40, 0, 0, targetLocation),
		}...)
	}

	times = append(times, []time.Time{
		time.Date(1997, time.September, 3, 9, 0, 0, 0, targetLocation),
		time.Date(1997, time.September, 3, 9, 20, 0, 0, targetLocation),
		time.Date(1997, time.September, 3, 9, 40, 0, 0, targetLocation),
	}...)

	RuleShouldMatchDates(t,
		s,
		times,
	)

	RuleShouldMatchDates(t,
		s2,
		times,
	)
}

func Test_BySecond(t *testing.T) {
	// There are no examples in the RFC that show BYSECOND
	// rules, but this should suffice for now.

	var s string = "DTSTART;TZID=America/New_York:19970902T090000\n" +
		"RRULE:FREQ=HOURLY;BYMINUTE=0,20,40;BYSECOND=10,20,30"

	RuleShouldMatchDates(t,
		s,
		[]time.Time{
			time.Date(1997, time.September, 2, 9, 0, 10, 0, targetLocation),
			time.Date(1997, time.September, 2, 9, 0, 20, 0, targetLocation),
			time.Date(1997, time.September, 2, 9, 0, 30, 0, targetLocation),
			time.Date(1997, time.September, 2, 9, 20, 10, 0, targetLocation),
			time.Date(1997, time.September, 2, 9, 20, 20, 0, targetLocation),
			time.Date(1997, time.September, 2, 9, 20, 30, 0, targetLocation),
			time.Date(1997, time.September, 2, 9, 40, 10, 0, targetLocation),
			time.Date(1997, time.September, 2, 9, 40, 20, 0, targetLocation),
			time.Date(1997, time.September, 2, 9, 40, 30, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 0, 10, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 0, 20, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 0, 30, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 20, 10, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 20, 20, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 20, 30, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 40, 10, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 40, 20, 0, targetLocation),
			time.Date(1997, time.September, 2, 10, 40, 30, 0, targetLocation),
		},
	)
}
