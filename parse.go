package rrule

import (
	"bytes"
	"time"
)

func Parse(rule string) (*RecurringRule, error) {
	recur_rule := RecurringRule{
		Interval:      1,           // default
		WorkWeekStart: time.Monday, // default
	}

	var chunk bytes.Buffer

	for _, c := range rule {
		if c == rune('\n') {
			err := recur_rule.handle_rule_chunk(chunk.String())
			if err != nil {
				return &recur_rule, err
			}
			chunk.Reset()
		} else {
			chunk.WriteRune(c)
		}
	}

	if chunk.Len() > 0 {
		err := recur_rule.handle_rule_chunk(chunk.String())
		if err != nil {
			return &recur_rule, err
		}
	}

	err := recur_rule.internal_parser()

	return &recur_rule, err
}
