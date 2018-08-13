package rrule

import (
	"time"
)

// Creating occurences of a recurrece rule is non-trivial because some
// fields allow positional or negative selection. In some cases possibilities
// are generated and the correct events are selected from them.
// In other cases (specifically, hours, minutes, seconds) those times
// are easier to generate.
type RecurrenceIterator struct {
	rule        *RecurringRule
	iterCounter int
	iterBuffer  []time.Time

	ReturnCounter int

	UserLimit int

	UserBefore time.Time
	UserAfter  time.Time

	UseUserBefore bool
	UseUserAfter  bool
}

func (ri *RecurrenceIterator) Limit(i int) *RecurrenceIterator {
	ri.UserLimit = i
	return ri
}

func (ri *RecurrenceIterator) Before(b time.Time) *RecurrenceIterator {
	ri.UserBefore = b
	ri.UseUserBefore = true
	return ri
}

func (ri *RecurrenceIterator) After(b time.Time) *RecurrenceIterator {
	ri.UserAfter = b
	ri.UseUserAfter = true
	return ri
}

func (ri *RecurrenceIterator) Between(a, b time.Time) *RecurrenceIterator {
	ri.After(a)
	ri.Before(b)
	return ri
}

func (ri *RecurrenceIterator) generateCandidates(root time.Time) []time.Time {
	var results []time.Time

	switch ri.rule.Frequency {
	case YEARLY:
		for m := time.January; m < 13; m += 1 {
			daysInMonth := time.Date(root.Year(), time.Month(m)+1, 0, 0, 0, 0, 0, time.UTC).Day()
			for d := 1; d <= daysInMonth; d += 1 {
				var target time.Time = time.Date(
					root.Year(),
					time.Month(m),
					d,
					root.Hour(),
					root.Minute(),
					root.Second(),
					0,
					root.Location(),
				)

				results = append(results, target)
			}
		}

	case MONTHLY:
		daysInMonth := time.Date(root.Year(), root.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		for d := 1; d <= daysInMonth; d += 1 {
			var target time.Time = time.Date(
				root.Year(),
				root.Month(),
				d,
				root.Hour(),
				root.Minute(),
				root.Second(),
				0,
				root.Location(),
			)

			results = append(results, target)
		}
	case WEEKLY:
		t := time.Date(root.Year(), root.Month(), root.Day(), root.Hour(), root.Minute(), root.Second(), 0, root.Location())

		// If the week start is non-standard
		// walk back until we find it.
		for t.Weekday() != ri.rule.WorkWeekStart {
			t = t.AddDate(0, 0, -1)
		}

		for d := 0; d < 7; d += 1 {
			results = append(results, t)
			t = t.AddDate(0, 0, 1)
		}

	case DAILY:
		fallthrough
	case HOURLY:
		fallthrough
	case MINUTELY:
		fallthrough
	case SECONDLY:
		results = append(results, ri.generateTimesForCandidates(root)...)
	default:
		panic("Unknown frequency")
	}

	return results
}

func (ri *RecurrenceIterator) generateTimesForCandidates(root time.Time) []time.Time {
	var results []time.Time

	var ByHourItems []int16 = []int16{int16(root.Hour())}
	var ByMinuteItems []int16 = []int16{int16(root.Minute())}
	var BySecondItems []int16 = []int16{int16(root.Second())}

	if len(ri.rule.ByHour) > 0 && (ri.rule.Frequency <= DAILY) {
		ByHourItems = ri.rule.ByHour
	}

	if len(ri.rule.ByMinute) > 0 && (ri.rule.Frequency <= MINUTELY) {
		ByMinuteItems = ri.rule.ByMinute
	}

	if len(ri.rule.BySecond) > 0 && (ri.rule.Frequency <= SECONDLY) {
		BySecondItems = ri.rule.BySecond
	}

	for _, h := range ByHourItems {
		for _, m := range ByMinuteItems {
			for _, s := range BySecondItems {
				results = append(
					results,
					time.Date(
						root.Year(),
						root.Month(),
						root.Day(),
						int(h),
						int(m),
						int(s),
						0,
						root.Location(),
					),
				)
			}
		}
	}

	return results
}

func (ri *RecurrenceIterator) Step(t *time.Time) bool {
	if ri.rule.Count > 0 && ri.ReturnCounter >= ri.rule.Count {
		return false
	}

	var shortcircuit_finish bool = false

	for len(ri.iterBuffer) == 0 && shortcircuit_finish == false {
		// If we've made it here, we need to generate more results
		next_base := getNextDateByFreqAndInterval(
			ri.rule.DtStart,
			ri.rule.Frequency,
			ri.rule.Interval*ri.iterCounter,
		)

		candidates := ri.generateCandidates(
			next_base,
		)

		candidate_count := len(candidates)

		final_set := []time.Time{}

		for cindex, d := range candidates {
			var matches []bool

			if len(ri.rule.ByYearDay) > 0 {
				var match bool = false
				for _, value := range ri.rule.ByYearDay {
					if value > 0 {
						if d.YearDay() == int(value) {
							match = true
						}
					} else if value < 0 {
						if ((candidate_count - cindex) + int(value)) == 0 {
							match = true
						}
					}
				}
				matches = append(matches, match)
			}

			if len(ri.rule.ByMonthDay) > 0 {
				var match bool = false
				for _, value := range ri.rule.ByMonthDay {
					if value > 0 {
						if d.Day() == int(value) {
							match = true
						}
					} else if value < 0 {
						if ((candidate_count - cindex) + int(value)) == 0 {
							match = true
						}
					}
				}
				matches = append(matches, match)
			}

			if len(ri.rule.ByMonth) > 0 {
				var match bool = false
				for _, value := range ri.rule.ByMonth {
					if value > 0 {
						if d.Month() == time.Month(value) {
							match = true
						}
					} else if value < 0 {
						if ((candidate_count - cindex) + int(value)) == 0 {
							match = true
						}
					}
				}
				matches = append(matches, match)
			}

			if len(ri.rule.ByWeekNo) > 0 {
				var match bool = false
				for _, value := range ri.rule.ByWeekNo {
					_, weeknum := d.ISOWeek()
					if value > 0 {
						if weeknum == int(value) {
							match = true
						}
					} else if value < 0 {
						if ((candidate_count - cindex) + int(value)) == 0 {
							match = true
						}
					}
				}

				matches = append(matches, match)
			}

			// We can't actually check the ByDay until we've evaluated
			// all of the candidates. But we'll filter here, to make
			// sure that there is less work in the future.

			if len(ri.rule.ByDay) > 0 {
				var match bool = false
				for _, value := range ri.rule.ByDay {
					if value.Weekday == d.Weekday() {
						if value.Offset == 0 {
							// match every occurence of a DOW
							match = true
						} else if value.Offset > 0 {
							// the first or Nth occurence
							if (cindex/7)+1 == value.Offset {
								match = true
							}
						} else if value.Offset < 0 {
							// the last or Nth last occurence
							cc_offset := ((candidate_count - cindex) - 1) / 7
							if cc_offset+(value.Offset+1) == 0 {
								match = true
							}
						}
					}
				}
				matches = append(matches, match)
			}

			// ByHour, ByMinute and BySecond don't support negative/relative
			// indexes, so the comparison is significanly simpler.

			if len(ri.rule.ByHour) > 0 {
				var match bool = false
				for _, value := range ri.rule.ByHour {
					if int(value) == d.Hour() {
						match = true
					}
				}

				matches = append(matches, match)
			}

			if len(ri.rule.ByMinute) > 0 {
				var match bool = false
				for _, value := range ri.rule.ByMinute {
					if int(value) == d.Minute() {
						match = true
					}
				}

				matches = append(matches, match)
			}

			if len(ri.rule.BySecond) > 0 {
				var match bool = false
				for _, value := range ri.rule.BySecond {
					if int(value) == d.Second() {
						match = true
					}
				}

				matches = append(matches, match)
			}

			// For a better understanding of some of the interesting
			// if | else if  cases below, take a look at
			// https://tools.ietf.org/html/rfc5545#page-44
			// there are some special cases to be aware of that
			// are not immediately obvious.

			if len(ri.rule.ByDay) == 0 && ri.rule.Frequency == WEEKLY {
				if d.Weekday() == ri.rule.DtStart.Weekday() {
					matches = append(matches, true)
				}
			}

			if len(ri.rule.ByHour) == 0 && ri.rule.Frequency == DAILY {
				if d.Hour() == ri.rule.DtStart.Hour() {
					matches = append(matches, true)
				}
			}

			if len(ri.rule.ByMinute) == 0 && ri.rule.Frequency == HOURLY {
				if d.Minute() == ri.rule.DtStart.Minute() {
					matches = append(matches, true)
				}
			}

			if len(ri.rule.BySecond) == 0 && ri.rule.Frequency == MINUTELY {
				if d.Second() == ri.rule.DtStart.Second() {
					matches = append(matches, true)
				}
			}

			if len(ri.rule.ByYearDay) == 0 && len(ri.rule.ByMonthDay) == 0 &&
				len(ri.rule.ByDay) == 0 && ri.rule.Frequency == YEARLY {

				var match bool = false
				if d.Day() == ri.rule.DtStart.Day() {
					match = true
				}
				matches = append(matches, match)
			}

			var final_match bool = true

			if len(matches) == 0 {
				final_match = false
			} else {
				for _, m := range matches {
					final_match = final_match && m
				}
			}

			if final_match {
				final_set = append(final_set, d)
			}
		}

		// Once we've generated the set of candidates that pass
		// the rules defined by the user, bysetpos can select from
		// that set to determine the final set.

		final_candidate_count := len(final_set)

		for cindex, d := range final_set {
			var matches []bool

			if len(ri.rule.BySetPos) > 0 {
				var match bool = false
				for _, value := range ri.rule.BySetPos {
					if value > 0 {
						if cindex == int(value-1) {
							match = true
						}
					} else if value < 0 {
						fcc := (final_candidate_count - cindex) + int(value)
						if fcc == 0 {
							match = true
						}
					}
				}

				matches = append(matches, match)
			}

			for _, exDate := range ri.rule.ExceptionsToRule {
				if exDate.Equal(d) {
					matches = append(matches, false)
				}
			}

			if ri.UseUserBefore {
				if d.Before(ri.UserBefore) {
					matches = append(matches, true)
				} else {
					matches = append(matches, false)
					// Since we can only generate events after this one
					// we are done.
					shortcircuit_finish = true
				}
			}

			if ri.UseUserAfter {
				if d.After(ri.UserAfter) {
					matches = append(matches, true)
				} else {
					matches = append(matches, false)
				}
			}

			var final_match bool = true

			if len(matches) > 0 {
				for _, m := range matches {
					final_match = final_match && m
				}
			}

			if final_match && !d.Before(ri.rule.DtStart) {
				ri.iterBuffer = append(ri.iterBuffer, d)
			}
		}

		ri.iterCounter += 1
	}

	if len(ri.iterBuffer) > 0 {
		*t = ri.iterBuffer[0]
		ri.iterBuffer = ri.iterBuffer[1:]

		if !ri.rule.Until.Equal(EmptyTime) && (*t).After(ri.rule.Until) == true {
			return false
		}

		ri.ReturnCounter += 1

		if ri.UserLimit > 0 && ri.ReturnCounter > ri.UserLimit {
			return false
		}
		return true
	}

	return !shortcircuit_finish
}

func getNextDateByFreqAndInterval(last time.Time, freq FrequencyValue, interval int) time.Time {
	switch freq {
	case YEARLY:
		// In order to avoid the double month increase on shorter months.
		last = time.Date(last.Year(), last.Month(), 1,
			last.Hour(), last.Minute(), last.Second(),
			0, last.Location())

		last = last.AddDate(1*interval, 0, 0)
	case MONTHLY:
		// In order to avoid the double month increase on shorter months.
		last = time.Date(last.Year(), last.Month(), 1,
			last.Hour(), last.Minute(), last.Second(),
			0, last.Location())

		last = last.AddDate(0, 1*interval, 0)
	case WEEKLY:
		last = last.AddDate(0, 0, 7*interval)
	case DAILY:
		last = last.AddDate(0, 0, 1*interval)
	case HOURLY:
		last = last.Add(time.Hour * time.Duration(interval))
	case MINUTELY:
		last = last.Add(time.Minute * time.Duration(interval))
	case SECONDLY:
		last = last.Add(time.Second * time.Duration(interval))
	}
	return last
}
