package rrule

import (
	"log"
	"testing"
	"time"
)

func Test_DateTime_Parser_UTC(t *testing.T) {
	var value = "19980119T070000Z"

	dt, err := ParseDateTime(value, time.UTC)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	match := time.Date(1998, time.January, 19, 7, 0, 0, 0, time.UTC)

	if !dt.Equal(match) {
		log.Printf("Date doesnt match %s", dt)
		t.Fail()
	}
}

func Test_DateTime_Parser_Local(t *testing.T) {
	var value = "19980119T070000"

	dt, err := ParseDateTime(value, time.Local)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	match := time.Date(1998, time.January, 19, 7, 0, 0, 0, time.Local)

	if !dt.Equal(match) {
		log.Printf("Date doesnt match %s", dt)
		t.Fail()
	}
}

func Test_Date_Parser_Local(t *testing.T) {
	var value = "19970714"

	dt, err := ParseDateTime(value, time.UTC)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	match := time.Date(1997, time.July, 14, 0, 0, 0, 0, time.UTC)

	if !dt.Equal(match) {
		log.Printf("Date doesnt match %s", dt)
		t.Fail()
	}
}

func Test_Parse_Duration_Complex(t *testing.T) {
	var value string = "P15DT5H0M20S"

	d, err := ParseDuration(value)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	dur := (15 * 24 * time.Hour) + (5 * time.Hour) + (20 * time.Second)

	if d != time.Duration(dur) {
		t.Log("Failed to parse duration correctly")
		t.Fail()
	}

}

func Test_Parse_Duration_Simple(t *testing.T) {
	var value string = "P7W"

	d, err := ParseDuration(value)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	dur := (7 * 7 * 24 * time.Hour)

	if d != time.Duration(dur) {
		t.Log("Failed to parse duration correctly")
		t.Fail()
	}

}

func Test_DateTimeToString(t *testing.T) {
	var value = "19980119T070000Z"
	var valueTime = time.Date(1998, time.January, 19, 7, 0, 0, 0, time.UTC)

	if value != DateTimeToString(valueTime) {
		t.Log("Datetime to string failed for UTC")
		t.Fail()
	}

	var value2 = "19980119T070000"
	var valueTime2 = time.Date(1998, time.January, 19, 7, 0, 0, 0, time.Local)

	if value2 != DateTimeToString(valueTime2) {
		t.Log("Datetime to string failed for local")
		t.Fail()
	}
}
