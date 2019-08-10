package eventexpander

import (
	"strings"
	"testing"
	"time"

	"github.com/graham/rrule"
	google_calendar "google.golang.org/api/calendar/v3"
)

var solo_compressedTestData []*google_calendar.Event
var solo_singleTestData []*google_calendar.Event

func init() {
	solo_compressedTestData = LoadTestData("solo_event_data/compressed_test_data.json")
	solo_singleTestData = LoadTestData("solo_event_data/single_event_test_data.json")
}

func TestSoloExport(t *testing.T) {
	source_event := solo_compressedTestData[0]

	if len(source_event.Recurrence) == 0 {
		t.Fatalf("Source event has no recurrence.")
	}

	rule, _ := rrule.Parse(strings.Join(source_event.Recurrence, "\n"))
	timeZone, _ := time.LoadLocation(source_event.Start.TimeZone)
	rule.DtStart, _ = time.Parse(time.RFC3339, source_event.Start.DateTime)
	rule.DtStart = rule.DtStart.In(timeZone)

	iter := rule.Iterator().
		After(time.Date(2019, time.January, 1, 1, 0, 0, 0, timeZone)).
		Before(time.Date(2019, time.July, 1, 1, 0, 0, 0, timeZone))

	var start time.Time

	var index int64 = 0

	for iter.Step(&start) {
		emulated_start := start.In(timeZone)
		reference_event := solo_singleTestData[index]
		index += 1

		real_start, _ := time.Parse(time.RFC3339, reference_event.Start.DateTime)

		if !real_start.Equal(emulated_start) {
			t.Fatalf("%s does not match %s", real_start, emulated_start)
		}

		emulated_event := CreateInstanceFromSource(source_event, emulated_start)

		if emulated_event.Id != reference_event.Id {
			t.Fatalf("Ids do not match: %s != %s", emulated_event.Id, reference_event.Id)
		}

		if emulated_event.Start.DateTime != reference_event.Start.DateTime {
			t.Fatalf("Start Times do not match: %s != %s", emulated_event.Start.DateTime, reference_event.Start.DateTime)
		}

		if emulated_event.End.DateTime != reference_event.End.DateTime {
			t.Fatalf("End Times do not match: %s != %s", emulated_event.End.DateTime, reference_event.End.DateTime)
		}

		if emulated_event.RecurringEventId != reference_event.RecurringEventId {
			t.Fatalf("recurring event ids do not match.")
		}
	}
}
