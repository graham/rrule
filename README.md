# rrule

RRule is a go library for expanding iCalendar recurrence into actual instances. It provides parsing for recurrence rule text, modification and re-exporting to a string recurrence rule.

I was able to find libraries that did each of these features individually, but not one that did all three. Lucikly this RFC https://tools.ietf.org/html/rfc5545 makes it very clear how rrules work (along with a ton of examples).

The best way to understand how to use the library would be to look at the tests, I've implemented every example in the RFC as well as some others that I created from real world problems I've run into.

GoDoc: https://godoc.org/github.com/graham/rrule

## Parsing from string
```
	rule_checker, err := rrule.Parse(
		"DTSTART;TZID=America/New_York:19970902T090000\n"+
			"RRULE:FREQ=DAILY;COUNT=10",
	)
```

## Recurrence Rule to String
```
	fmt.Println( rule.String() )
```

## Iterating through instances

```
	rule_checker, err := rrule.Parse(
		"DTSTART;TZID=America/New_York:19970902T090000\n"+
			"RRULE:FREQ=DAILY;COUNT=10",
	)

	iter := rule.Iterator()

	for iter.Step(&event) {
		fmt.Println(event)
	}
```

## Questions and Contributions
If you find a bug, please file an issue, or propose a change, I'd be happy to include changes if you find a case I haven't covered.

This was a great opportunity for me to learn about the iCalendar standard, but ergonomically the library could be easier to use (expecially more examples). Recommendations welcome.

In terms of events themselves, I tried to ensure the core rrule library didn't have any dependencies to in the wild implementations (Google Calendar for example).

I have, however, started an example application under `eventexpander` that should correctly output event instances if given a Google Calendar event json object. As of 2019-08-10 I'm still working on this, but if you're interested in an actual implementation I'd recommend taking a look at that code.

Have a great day :)
