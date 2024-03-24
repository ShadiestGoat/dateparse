package dateparse

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var minTime = time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)

type ComboParseErr struct {
	Value  string
	Errors []error
}

func (err ComboParseErr) Error() string {
	if len(err.Errors) == 0 {
		return fmt.Sprintf("Failed to parse '%v': unknown (and possibly no) errors", err.Value)
	}

	errors := ""
	for _, e := range err.Errors {
		errors += "\t\"" + e.Error() + "\"\n"
	}
	return fmt.Sprintf("Failed to parse '%v':\n%v", err.Value, errors[:len(errors)-1])
}

// Parses a str using multiple layouts, returning the successful layout.
func parse(str string, layouts []string, tz *time.Location) (time.Time, string, error) {
	errors := []error{}

	for _, l := range layouts {
		t, err := time.ParseInLocation(l, str, tz)
		if err == nil {
			return t, l, nil
		}
		errors = append(errors, err)
	}

	return time.Time{}, "", ComboParseErr{
		Value:  str,
		Errors: errors,
	}
}

type ParseError struct {
	Layout, Value string
}

func (err ParseError) Error() string {
	return "Parse Error: cannot parse '" + err.Value + "' using '" + err.Layout + "'"
}

// Parses unix time in different formats:
// 10: seconds
// 13: milliseconds
// 16: microseconds
// 19: nanoseconds
func ParseUnix(inp string) (time.Time, error) {
	l := len(inp)
	pErr := ParseError{
		// Real layout, trust me
		Layout: "unix timestamp",
		Value:  inp,
	}

	if l < 10 || l > 19 || (l-1)%3 != 0 {
		return time.Time{}, pErr
	}

	v, err := strconv.ParseInt(inp, 10, 64)
	if err != nil {
		return time.Time{}, pErr
	}

	var t time.Time

	switch l {
	case 10:
		t = time.Unix(v, 0)
	case 13:
		t = time.UnixMilli(v)
	case 16:
		t = time.UnixMicro(v)
	case 19:
		nsec := v % 1000
		t = time.Unix(v-nsec, nsec)
	}

	return t, nil
}

func isTime(inp string) bool {
	if strings.Contains(inp, ":") {
		return true
	}

	if len(inp) <= 2 {
		return false
	}

	suf := strings.ToLower(inp[len(inp)-2:])

	if len(inp) > 4 || (suf != "pm" && suf != "am") {
		return false
	}

	inp = inp[:len(inp)-2]

	for _, r := range inp {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

// Today, at 00:00:00
func Today(tz *time.Location) time.Time {
	now := time.Now()

	if tz != nil {
		now = now.In(tz)
	} else {
		tz = time.UTC
	}

	y, m, d := now.Date()

	return time.Date(y, m, d, 0, 0, 0, 0, tz)
}
