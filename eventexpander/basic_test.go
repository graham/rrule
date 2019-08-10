package eventexpander

import (
	"bufio"
	"encoding/json"
	"os"
	"sort"

	google_calendar "google.golang.org/api/calendar/v3"
)

func LoadTestData(filename string) []*google_calendar.Event {
	var openErr error
	var f *os.File

	var result []*google_calendar.Event

	f, openErr = os.Open(filename)

	if openErr != nil {
		panic(openErr)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := scanner.Bytes()
		var r google_calendar.Event
		err := json.Unmarshal(text, &r)

		if err != nil {
			panic(err)
		}

		result = append(result, &r)
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })

	return result
}
