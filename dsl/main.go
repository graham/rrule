package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/graham/rrule"
)

func main() {
	var rule_string string = "RRULE:FREQ=WEEKLY;INTERVAL=2;COUNT=3"
	var last_time_it_happened = time.Date(2018, 1, 19, 10, 0, 0, 0, time.Local)
	var rr *rrule.RecurringRule
	var err error

	rr, err = rrule.Parse(rule_string)

	if err != nil {
		panic(err)
	}

	rr.DtStart = last_time_it_happened

	fmt.Println(rr)

	iter := rr.Iterator()

	var occurence time.Time

	enc := json.NewEncoder(os.Stdout)

	for iter.Step(&occurence) {
		enc.Encode(occurence)
	}
}
