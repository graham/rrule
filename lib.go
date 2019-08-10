// RRule is a go package for parsing and editing iCalendar Recurrence Rules. It also
// provides the code to generate occurences of those events (given additiona details).
//
// Primary access to rules is done via the parse command.
//
//   rule, err := rrule.Parse("RRULE:FREQ=DAILY;COUNT=10")
//
// Access to the occurences of a rule is via the Iterator struct.
//
//   iter := rule.Iterator()
//
//   for iter.Step(&event) {
//       ... process your event here.
//   }
//
// Step will return false when you more results. Some events
// recur forever, so be careful to ensure your loop will end.
//
// Iterators have a DtStart field, (which can also be parsed,
// via the parse command. Setting the DtStart manually is allowed
// and often required if that data is stored separately
// (Google Calendar for example).
//
// Iterators have the helper methods Limit, Before, After and Between to
// allow for easier filtering.
//
// In addition to the code for parsing and iterating through events,
// the library includes all of the example dates included in the icalendar
// rfc (5545) (including some of the erratta).
package rrule

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// https://tools.ietf.org/html/rfc5545#page-122
type FrequencyValue uint8

// EmptyTime is an uninitialized time value, used to compare to other
// uninitialized time values.
var EmptyTime time.Time

// Constants
const (
	YEARLY FrequencyValue = iota
	MONTHLY
	WEEKLY
	DAILY
	HOURLY
	MINUTELY
	SECONDLY
)

func (fv FrequencyValue) String() string {
	switch fv {
	case YEARLY:
		return "YEARLY"
	case MONTHLY:
		return "MONTHLY"
	case WEEKLY:
		return "WEEKLY"
	case DAILY:
		return "DAILY"
	case HOURLY:
		return "HOURLY"
	case MINUTELY:
		return "MINUTELY"
	case SECONDLY:
		return "SECONDLY"
	}
	panic("Unknown Frequency")
}

func weekdayFromShort(s string) time.Weekday {
	switch s {
	case "SU":
		return time.Sunday
	case "MO":
		return time.Monday
	case "TU":
		return time.Tuesday
	case "WE":
		return time.Wednesday
	case "TH":
		return time.Thursday
	case "FR":
		return time.Friday
	case "SA":
		return time.Saturday
	}
	panic("Unknown Day")
}

func shortFromWeekday(t time.Weekday) string {
	switch t {
	case time.Sunday:
		return "SU"
	case time.Monday:
		return "MO"
	case time.Tuesday:
		return "TU"
	case time.Wednesday:
		return "WE"
	case time.Thursday:
		return "TH"
	case time.Friday:
		return "FR"
	case time.Saturday:
		return "SA"
	}
	panic("Unknown day")
}

type ForDay struct {
	Weekday time.Weekday
	Offset  int
}

func (fd ForDay) String() string {
	s := shortFromWeekday(fd.Weekday)
	if fd.Offset == 0 {
		return s
	} else {
		return fmt.Sprintf("%d%s", fd.Offset, s)
	}
}

type RecurringRule struct {
	DtStart time.Time

	// https://tools.ietf.org/html/rfc5545#page-39
	Frequency     FrequencyValue
	Until         time.Time
	Count         int
	Interval      int
	BySecond      []int16  //    0 - 60
	ByMinute      []int16  //    0 - 59
	ByHour        []int16  //    0 - 23
	ByDay         []ForDay //  -53 - 53
	ByMonthDay    []int16  //  -31 - 31
	ByYearDay     []int16  // -366 - 366
	ByWeekNo      []int16  //  -53 - 53
	ByMonth       []int16  //    1 - 12
	BySetPos      []int16  // -366 - 366
	WorkWeekStart time.Weekday

	ExceptionsToRule []time.Time
}

func compareListsOfInt16(a, b []int16) bool {
	if len(a) != len(b) {
		return false
	}

	for index, _ := range a {
		if a[index] != b[index] {
			return false
		}
	}

	return true
}

func (r1 *RecurringRule) Equal(r2 *RecurringRule) bool {
	if !r1.DtStart.Equal(r2.DtStart) {
		return false
	}

	if r1.Frequency != r2.Frequency {
		return false
	}

	if !r1.Until.Equal(r2.Until) {
		return false
	}

	if r1.Count != r2.Count {
		return false
	}

	if r1.Interval != r2.Interval {
		return false
	}

	if compareListsOfInt16(r1.BySecond, r2.BySecond) == false {
		return false
	}

	if compareListsOfInt16(r1.ByMinute, r2.ByMinute) == false {
		return false
	}

	if compareListsOfInt16(r1.ByHour, r2.ByHour) == false {
		return false
	}

	if len(r1.ByDay) != len(r2.ByDay) {
		return false
	} else {
		for index, _ := range r1.ByDay {
			if r1.ByDay[index] != r2.ByDay[index] {
				return false
			}
		}
	}

	if compareListsOfInt16(r1.ByMonthDay, r2.ByMonthDay) == false {
		return false
	}

	if compareListsOfInt16(r1.ByYearDay, r2.ByYearDay) == false {
		return false
	}

	if compareListsOfInt16(r1.ByWeekNo, r2.ByWeekNo) == false {
		return false
	}

	if compareListsOfInt16(r1.BySetPos, r2.BySetPos) == false {
		return false
	}

	if r1.WorkWeekStart != r2.WorkWeekStart {
		return false
	}

	if len(r1.ExceptionsToRule) != len(r2.ExceptionsToRule) {
		return false
	}

	for index, exDate := range r1.ExceptionsToRule {
		if !exDate.Equal(r2.ExceptionsToRule[index]) {
			return false
		}
	}

	return true
}

func (rr *RecurringRule) internal_parser() error {
	if passed, reason := rr.passed_rfc_checks(); passed == false {
		return errors.New(reason)
	}

	return nil
}

func (rr *RecurringRule) passed_rfc_checks() (bool, string) {
	for _, item := range rr.BySecond {
		if item < 0 || item > 60 {
			return false, "BySecond Rule fail: 0 <= value <= 60"
		}
	}

	for _, item := range rr.ByMinute {
		if item < 0 || item > 59 {
			return false, "ByMinute Rule fail: 0 <= value <= 59"
		}
	}

	for _, item := range rr.ByHour {
		if item < 0 || item > 23 {
			return false, "ByHour Rule fail: 0 <= value <= 23"
		}
	}

	for _, item := range rr.ByDay {
		if item.Offset < -53 || item.Offset > 53 {
			return false, "ByDay Rule fail: -53 <= value <= 53"
		}
	}

	return true, ""
}

func (rr *RecurringRule) handle_part_freq(value string) error {
	switch value {
	case "SECONDLY":
		rr.Frequency = SECONDLY
	case "MINUTELY":
		rr.Frequency = MINUTELY
	case "HOURLY":
		rr.Frequency = HOURLY
	case "DAILY":
		rr.Frequency = DAILY
	case "WEEKLY":
		rr.Frequency = WEEKLY
	case "MONTHLY":
		rr.Frequency = MONTHLY
	case "YEARLY":
		rr.Frequency = YEARLY
	default:
		return errors.New(fmt.Sprintf("%s is not a valid FREQ", value))
	}

	return nil
}

func (rr *RecurringRule) handle_recur_rule_part(key string, value string) error {
	switch key {
	case "FREQ":
		if err := rr.handle_part_freq(value); err != nil {
			return err
		}
	case "UNTIL":
		dt, err := ParseDateTime(value, rr.DtStart.Location())
		if err != nil {
			return err
		}
		rr.Until = dt
	case "COUNT":
		iv, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		rr.Count = int(iv)
	case "INTERVAL":
		iv, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		rr.Interval = int(iv)
	case "BYSECOND":
		fallthrough
	case "BYMINUTE":
		fallthrough
	case "BYHOUR":
		fallthrough
	case "BYMONTHDAY":
		fallthrough
	case "BYYEARDAY":
		fallthrough
	case "BYWEEKNO":
		fallthrough
	case "BYMONTH":
		fallthrough
	case "BYSETPOS":
		var results []int16
		chunks := strings.Split(value, ",")
		for _, item := range chunks {
			iv, err := strconv.ParseInt(item, 10, 64)
			if err != nil {
				return err
			}
			results = append(results, int16(iv))
		}
		switch key {
		case "BYSECOND":
			rr.BySecond = results
		case "BYMINUTE":
			rr.ByMinute = results
		case "BYHOUR":
			rr.ByHour = results
		case "BYMONTHDAY":
			rr.ByMonthDay = results
		case "BYYEARDAY":
			rr.ByYearDay = results
		case "BYWEEKNO":
			rr.ByWeekNo = results
		case "BYMONTH":
			rr.ByMonth = results
		case "BYSETPOS":
			rr.BySetPos = results
		}
	case "WKST":
		rr.WorkWeekStart = weekdayFromShort(value)
	case "BYDAY":
		var results []ForDay
		chunks := strings.Split(value, ",")
		for _, item := range chunks {
			var offset_dir int = 1
			var offset int = 0
			var chunk string

			if item[0] == '-' {
				offset_dir = -1
				item = item[1:]
			}

			if unicode.IsDigit(rune(item[0])) && unicode.IsDigit(rune(item[1])) {
				chunk = item[0:2]
				item = item[2:]
			} else if unicode.IsDigit(rune(item[0])) {
				chunk = item[0:1]
				item = item[1:]
			}

			if len(chunk) > 0 {
				iv, err := strconv.ParseInt(chunk, 10, 64)
				if err != nil {
					return err
				}
				offset = int(iv)
			}

			dow := weekdayFromShort(item)
			results = append(results, ForDay{
				Weekday: dow,
				Offset:  offset * offset_dir,
			})
		}
		rr.ByDay = results
	}

	return nil
}

func (rr *RecurringRule) handle_rule_chunk(value string) error {
	if strings.HasPrefix(value, "RRULE:") {
		parts := strings.Split(value[6:], ";")
		for _, part := range parts {
			args := strings.SplitN(part, "=", 2)
			err := rr.handle_recur_rule_part(args[0], args[1])
			if err != nil {
				return err
			}
		}
	} else if strings.HasPrefix(value, "DTSTART;") {
		dt, err := ParseDateTimeChunks(value[8:])
		if err != nil {
			return err
		}
		rr.DtStart = dt[0]
	} else if strings.HasPrefix(value, "EXDATE;") {
		dts, err := ParseDateTimeChunks(value[7:])
		if err != nil {
			return err
		}
		rr.ExceptionsToRule = dts
	} else {
		panic(fmt.Sprintf("Unknown control: %s", value))
	}

	return nil
}

func (rr *RecurringRule) Iterator() *RecurrenceIterator {
	return &RecurrenceIterator{rule: rr}
}

func listOfIntsToCSV(values []int16) string {
	numbers := make([]string, 0)

	for _, iv := range values {
		numbers = append(numbers, fmt.Sprintf("%d", iv))
	}

	return strings.Join(numbers, ",")
}

func (rr *RecurringRule) ExdateString() string {
	if len(rr.ExceptionsToRule) > 0 {
		return fmt.Sprintf("EXDATE%s", DatesToList(rr.ExceptionsToRule))
	} else {
		return ""
	}
}

func (rr *RecurringRule) RecurString() string {
	var rules []string

	rules = append(rules, "RRULE:")

	// There is always a freq

	rules = append(rules,
		fmt.Sprintf("FREQ=%s", rr.Frequency.String()))

	if rr.Interval != 1 {
		rules = append(rules,
			fmt.Sprintf(";INTERVAL=%d", rr.Interval))
	}

	if !rr.Until.Equal(EmptyTime) {
		rules = append(rules,
			fmt.Sprintf(";UNTIL=%s", DateTimeToString(rr.Until)))
	}

	if len(rr.BySecond) > 0 {
		rules = append(rules,
			fmt.Sprintf(
				";BYSECOND=%s",
				listOfIntsToCSV(rr.BySecond)))
	}

	if len(rr.ByMinute) > 0 {
		rules = append(rules,
			fmt.Sprintf(
				";BYMINUTE=%s",
				listOfIntsToCSV(rr.ByMinute)))
	}

	if len(rr.ByHour) > 0 {
		rules = append(rules,
			fmt.Sprintf(
				";BYHOUR=%s",
				listOfIntsToCSV(rr.ByHour)))
	}

	if len(rr.ByMonth) > 0 {
		rules = append(rules,
			fmt.Sprintf(
				";BYMONTH=%s",
				listOfIntsToCSV(rr.ByMonth)))
	}

	if len(rr.ByWeekNo) > 0 {
		rules = append(rules,
			fmt.Sprintf(
				";BYWEEKNO=%s",
				listOfIntsToCSV(rr.ByWeekNo)))
	}

	if rr.Count > 0 {
		rules = append(rules, fmt.Sprintf(";COUNT=%d", rr.Count))
	}

	if rr.WorkWeekStart != time.Monday {
		rules = append(rules,
			fmt.Sprintf(";WKST=%s",
				shortFromWeekday(rr.WorkWeekStart)))
	}

	if len(rr.ByDay) > 0 {
		days_as_str := []string{}
		for _, fd := range rr.ByDay {
			days_as_str = append(days_as_str,
				fd.String())
		}

		rules = append(rules,
			fmt.Sprintf(
				";BYDAY=%s",
				strings.Join(days_as_str, ",")))
	}

	if len(rr.ByMonthDay) > 0 {
		rules = append(rules,
			fmt.Sprintf(
				";BYMONTHDAY=%s",
				listOfIntsToCSV(rr.ByMonthDay)))
	}

	if len(rr.ByYearDay) > 0 {
		rules = append(rules,
			fmt.Sprintf(
				";BYYEARDAY=%s",
				listOfIntsToCSV(rr.ByYearDay)))
	}

	if len(rr.BySetPos) > 0 {
		rules = append(rules,
			fmt.Sprintf(
				";BYSETPOS=%s",
				listOfIntsToCSV(rr.BySetPos)))
	}

	return strings.Join(rules, "")
}

func (rr *RecurringRule) String() string {
	result_string := []string{}

	if rr.DtStart.Equal(EmptyTime) == false {
		result_string = append(result_string,
			fmt.Sprintf(
				"DTSTART;TZID=%s:%s",
				rr.DtStart.Location().String(),
				DateTimeToString(rr.DtStart),
			))
	}

	if len(rr.ExceptionsToRule) > 0 {
		result_string = append(result_string, rr.ExdateString())
	}
	result_string = append(result_string, rr.RecurString())

	return strings.Join(result_string, "\n")
}
