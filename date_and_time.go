package rrule

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var BadFormatError func(string) error = func(key string) error { return errors.New(fmt.Sprintf("Invalid Format: %s", key)) }

func ParseDateTimeChunks(s string) ([]time.Time, error) {
	var dtDates []time.Time
	var dtLocation *time.Location = time.UTC
	var parts []string

	if strings.Contains(s, ":") {
		parts = strings.Split(s, ":")
		for _, item := range strings.Split(parts[0], ";") {
			args := strings.SplitN(item, "=", 2)
			if args[0] == "TZID" {
				loc, err := time.LoadLocation(args[1])
				if err != nil {
					return dtDates, err
				}
				dtLocation = loc
			} else {
				panic(fmt.Sprintf("Unknown key/value %s=%s", args[0], args[1]))
			}
		}
		parts = strings.Split(parts[1], ",")
	} else {
		parts = strings.Split(s, ",")
	}

	for _, item := range parts {
		dt, err := ParseDateTime(item, dtLocation)
		if err != nil {
			return dtDates, err
		}
		dtDates = append(dtDates, dt)
	}

	return dtDates, nil
}

func ParseDateTime(s string, InLocation *time.Location) (time.Time, error) {
	now := time.Now()
	s_len := len(s)

	if s_len < 8 {
		return now, BadFormatError("Too short")
	}

	year := s[0:4]
	month := s[4:6]
	day := s[6:8]

	y, err := strconv.ParseInt(year, 10, 32)
	if err != nil {
		return now, err
	}

	m, err := strconv.ParseInt(month, 10, 32)
	if err != nil {
		return now, err
	}

	d, err := strconv.ParseInt(day, 10, 32)
	if err != nil {
		return now, err
	}

	if len(s) == 8 {
		return time.Date(int(y), time.Month(m), int(d), 0, 0, 0, 0, time.UTC), nil
	}

	if s[8] != 'T' {
		return now, BadFormatError("No T for time")
	}

	if s_len < 15 || s_len > 16 {
		return now, BadFormatError("not correct lenght for datetime")
	}

	hours := s[9:11]
	minutes := s[11:13]
	seconds := s[13:15]
	is_utc := strings.HasSuffix(s, "Z")

	hrs, err := strconv.ParseInt(hours, 10, 32)
	if err != nil {
		return now, err
	}

	min, err := strconv.ParseInt(minutes, 10, 32)
	if err != nil {
		return now, err
	}

	sec, err := strconv.ParseInt(seconds, 10, 32)
	if err != nil {
		return now, err
	}

	var locale *time.Location = InLocation

	if is_utc {
		locale = time.UTC
	}

	return time.Date(
		int(y),
		time.Month(m),
		int(d),
		int(hrs),
		int(min),
		int(sec), 0, locale), nil
}

func ParseDuration(s string) (time.Duration, error) {
	no_duration := time.Duration(0)

	var pos_or_neg int64 = 1
	if s[0] == '-' {
		pos_or_neg = -1
		s = s[1:]
	}

	var chunk bytes.Buffer

	if s[0] != 'P' {
		return no_duration, BadFormatError(s)
	}

	s = s[1:]

	var durationSeconds int64 = 0

	for _, c := range s {
		if unicode.IsDigit(c) {
			chunk.WriteRune(c)
		} else if c == 'T' {
			// pass
		} else {
			value, err := strconv.ParseInt(chunk.String(), 10, 64)
			if err != nil {
				return no_duration, BadFormatError(fmt.Sprintf("%d", value))
			}
			switch c {
			case 'W':
				durationSeconds += (value * 60 * 60 * 24 * 7)
			case 'D':
				durationSeconds += (value * 60 * 60 * 24)
			case 'H':
				durationSeconds += (value * 60 * 60)
			case 'M':
				durationSeconds += (value * 60)
			case 'S':
				durationSeconds += (value)
			}
			chunk.Reset()
		}
	}

	durationSeconds = durationSeconds * pos_or_neg

	return time.Duration(durationSeconds) * time.Second, nil
}

func DateTimeToString(t time.Time) string {
	if t.Location() == time.UTC {
		return fmt.Sprintf("%04d%02d%02dT%02d%02d%02dZ",
			t.Year(),
			int(t.Month()),
			t.Day(),
			t.Hour(),
			t.Minute(),
			t.Second(),
		)
	} else {
		return fmt.Sprintf("%04d%02d%02dT%02d%02d%02d",
			t.Year(),
			int(t.Month()),
			t.Day(),
			t.Hour(),
			t.Minute(),
			t.Second(),
		)
	}
}

func DatesToList(times []time.Time) string {
	if len(times) == 0 {
		return ""
	}

	var loc = times[0].Location()
	var correctedTimes []time.Time = []time.Time{times[0]}

	// All times must be in the same timezone.
	for _, t2 := range times[1:] {
		correctedTimes = append(correctedTimes, t2.In(loc))
	}

	var correctedTimeStrings []string

	for _, ct := range correctedTimes {
		correctedTimeStrings = append(correctedTimeStrings,
			DateTimeToString(ct))
	}

	if loc == time.UTC {
		return fmt.Sprintf(
			":%s",
			strings.Join(correctedTimeStrings, ","),
		)
	} else {
		return fmt.Sprintf(
			";TZID=%s:%s",
			loc.String(),
			strings.Join(correctedTimeStrings, ","),
		)

	}
}
