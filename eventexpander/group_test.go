package eventexpander

import (
	"strings"
	"testing"
	"time"

	"github.com/graham/rrule"
	google_calendar "google.golang.org/api/calendar/v3"
)

var group_compressedTestData []*google_calendar.Event
var group_singleTestData []*google_calendar.Event

func init() {
	group_compressedTestData = LoadTestData("group_event_data/compressed_test_data.json")
	group_singleTestData = LoadTestData("group_event_data/single_event_test_data.json")
}

func TestGroupExport(t *testing.T) {
	source_event := group_compressedTestData[0]

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

	var exceptionEvents map[string]*google_calendar.Event = make(map[string]*google_calendar.Event, 0)
	for _, i := range group_compressedTestData[1:] {
		exceptionEvents[i.Id] = i
	}

	var index int = 0

	for iter.Step(&start) {
		emulated_event_id := GenerateEventId(source_event, start)
		var emulated_event *google_calendar.Event

		if ex, found := exceptionEvents[emulated_event_id]; found == true {
			emulated_event = ex
		} else {
			emulated_event = CreateInstanceFromSource(source_event, start)
		}

		reference_event := group_singleTestData[index]

		if emulated_event.Id != reference_event.Id {
			t.Fatalf("Ids do not match: %s != %s", emulated_event.Id, reference_event.Id)
		}

		if emulated_event.Start.DateTime != reference_event.Start.DateTime {
			t.Fatalf("Start Times do not match: %s != %s", emulated_event.Start.DateTime, reference_event.Start.DateTime)
		}

		if emulated_event.End.DateTime != reference_event.End.DateTime {
			t.Fatalf("End Times do not match: %s != %s", emulated_event.End.DateTime, reference_event.End.DateTime)
		}

		index += 1
	}
}
