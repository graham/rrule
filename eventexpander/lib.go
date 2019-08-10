package eventexpander

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	google_calendar "google.golang.org/api/calendar/v3"
)

func timeToEventUid(start time.Time) string {
	return start.In(time.UTC).Format("20060102T150405")
}

func GenerateRootEventId(source *google_calendar.Event, start time.Time) string {
	var id string = source.Id
	if strings.Contains(id, "_") {
		sp := strings.Split(id, "_")
		id = sp[0]
	}

	return fmt.Sprintf("%s_R%s", id, timeToEventUid(start))
}
func GenerateEventId(source *google_calendar.Event, start time.Time) string {
	var id string = source.Id
	if strings.Contains(id, "_") {
		sp := strings.Split(id, "_")
		id = sp[0]
	}

	return fmt.Sprintf("%s_%sZ", id, timeToEventUid(start))
}

func CreateInstanceFromSource(source *google_calendar.Event, start time.Time) *google_calendar.Event {
	eventStart, _ := time.Parse(time.RFC3339, source.Start.DateTime)
	eventEnd, _ := time.Parse(time.RFC3339, source.End.DateTime)

	durationMinutes := int(eventEnd.Sub(eventStart).Minutes())

	var resultEvent google_calendar.Event

	copier.Copy(&resultEvent, source)

	resultEvent.Recurrence = []string{}

	resultEvent.Start = &google_calendar.EventDateTime{
		DateTime: start.Format(time.RFC3339),
		TimeZone: source.Start.TimeZone,
	}
	resultEvent.End = &google_calendar.EventDateTime{
		DateTime: start.Add(time.Duration(durationMinutes) * time.Minute).Format(time.RFC3339),
		TimeZone: source.Start.TimeZone,
	}

	// Set id to normal id with _20190101T150000Z after.

	resultEvent.RecurringEventId = GenerateRootEventId(source, eventStart)
	resultEvent.Id = GenerateEventId(&resultEvent, start)

	return &resultEvent
}
